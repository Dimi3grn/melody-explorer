package api

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/yourusername/melody-explorer/internal/models"
	"github.com/yourusername/melody-explorer/internal/spotify"
	"github.com/yourusername/melody-explorer/internal/storage"
)

// Server représente le serveur API
type Server struct {
	Router           *mux.Router
	SpotifyAuth      *spotify.Auth
	SpotifyClient    *spotify.Client
	FavoritesStorage *storage.FavoritesStorage
	TemplatesDir     string
	StaticDir        string
	templates        map[string]*template.Template
}

// NewServer crée une nouvelle instance de serveur
func NewServer(templatesDir, staticDir, dataDir string) (*Server, error) {
	// Créer le routeur
	router := mux.NewRouter()

	// Créer l'authentification Spotify
	auth, err := spotify.NewAuth()
	if err != nil {
		return nil, fmt.Errorf("échec lors de la création de l'authentification Spotify: %w", err)
	}

	// Créer le client Spotify
	client := spotify.NewClient(auth)

	// Créer le stockage des favoris
	favoritesStorage, err := storage.NewFavoritesStorage(dataDir)
	if err != nil {
		return nil, fmt.Errorf("échec lors de la création du stockage des favoris: %w", err)
	}

	// Créer le serveur
	server := &Server{
		Router:           router,
		SpotifyAuth:      auth,
		SpotifyClient:    client,
		FavoritesStorage: favoritesStorage,
		TemplatesDir:     templatesDir,
		StaticDir:        staticDir,
		templates:        make(map[string]*template.Template),
	}

	// Analyser les templates
	if err := server.parseTemplates(); err != nil {
		return nil, fmt.Errorf("échec lors de l'analyse des templates: %w", err)
	}

	// Initialiser les routes
	server.initializeRoutes()

	return server, nil
}

// parseTemplates analyse tous les templates
func (s *Server) parseTemplates() error {
	// Créer une carte de fonctions pour les templates
	// Créer une carte de fonctions pour les templates
	funcMap := template.FuncMap{
		"formatDate": func(t time.Time) string {
			return t.Format("January 2, 2006")
		},
		"add": func(a, b int) int {
			return a + b
		},
		"formatDuration": func(ms int) string {
			duration := time.Duration(ms) * time.Millisecond
			minutes := int(duration.Minutes())
			seconds := int(duration.Seconds()) % 60
			return fmt.Sprintf("%d:%02d", minutes, seconds)
		},
		"div": func(a, b int) int {
			return a / b
		},
		"mod": func(a, b int) int {
			return a % b
		},
		"contains": func(slice []string, str string) bool {
			for _, s := range slice {
				if s == str {
					return true
				}
			}
			return false
		},
		// Ajouter ici les nouveaux assistants de pagination
		"buildPaginationURL": func(baseURL, queryString string, page int, additionalParams map[string]string) string {
			u, err := url.Parse(baseURL)
			if err != nil {
				return "#"
			}

			q := u.Query()
			if queryString != "" {
				q.Set("q", queryString)
			}
			q.Set("page", strconv.Itoa(page))

			// Ajouter des paramètres supplémentaires
			for key, value := range additionalParams {
				if value != "" {
					q.Set(key, value)
				}
			}

			u.RawQuery = q.Encode()
			return u.String()
		},
		"formatPageCount": func(currentPage, totalPages int) string {
			if totalPages > 100 {
				return fmt.Sprintf("Page %d sur plusieurs", currentPage)
			}
			return fmt.Sprintf("Page %d sur %d", currentPage, totalPages)
		},
	}

	layouts, err := filepath.Glob(filepath.Join(s.TemplatesDir, "layout", "*.html"))
	if err != nil {
		return err
	}

	includes, err := filepath.Glob(filepath.Join(s.TemplatesDir, "partials", "*.html"))
	if err != nil {
		return err
	}

	// Générer des templates pour chaque page
	pages, err := filepath.Glob(filepath.Join(s.TemplatesDir, "pages", "*.html"))
	if err != nil {
		return err
	}

	for _, page := range pages {
		files := append(layouts, page)
		files = append(files, includes...)

		name := filepath.Base(page)
		tmpl, err := template.New(name).Funcs(funcMap).ParseFiles(files...)
		if err != nil {
			return err
		}

		s.templates[name] = tmpl
	}

	return nil
}

// renderTemplate rend un template avec les données fournies
func (s *Server) renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	tmpl, ok := s.templates[name]
	if !ok {
		// Au lieu d'appeler http.Error qui écrit un en-tête et un corps
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Template introuvable"))
		return
	}

	// Définir le type de contenu
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Exécuter le template directement sur l'écrivain
	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		// Journaliser l'erreur mais ne pas essayer d'écrire une nouvelle réponse si nous avons déjà commencé
		log.Printf("Erreur lors de l'exécution du template: %v", err)
	}
}

// PageData représente les données transmises aux templates
type PageData struct {
	Title       string
	IsLoggedIn  bool
	CurrentPage string
	Data        interface{}
	Query       string
	Filters     map[string]string
	Pagination  *PaginationData
	Error       string
}

// PaginationData représente les informations de pagination
type PaginationData struct {
	CurrentPage int
	TotalPages  int
	TotalItems  int
	Limit       int
	HasPrev     bool
	HasNext     bool
	PrevPage    int
	NextPage    int
}

// HomeHandler gère la page d'accueil
func (s *Server) HomeHandler(w http.ResponseWriter, r *http.Request) {
	isLoggedIn := s.SpotifyAuth.IsTokenValid()

	// Valeurs par défaut pour les données
	var recentTracks []spotify.Track

	// Liste fixe de genres populaires (au lieu d'essayer de les obtenir depuis l'API)
	genres := []string{
		"pop", "rock", "hip-hop", "jazz", "electronic",
		"classical", "country", "reggae", "blues",
		"metal", "indie", "folk", "r&b", "latin",
		"alternative", "dance", "ambient", "punk",
	}

	if isLoggedIn {
		// Obtenir des recommandations (qui utilisent maintenant les nouvelles sorties)
		var err error
		recentTracks, err = s.SpotifyClient.GetRecommendations(nil, 10)
		if err != nil {
			log.Printf("Erreur lors de l'obtention des recommandations: %v\n", err)
			// Nous allons simplement continuer avec une liste de pistes vide
		}
	}

	// Préparer les données de la page
	data := PageData{
		Title:       "MelodyExplorer - Découvrir de la Musique",
		IsLoggedIn:  isLoggedIn,
		CurrentPage: "home",
		Data: map[string]interface{}{
			"Tracks": recentTracks,
			"Genres": genres,
		},
	}

	// Rendre le template
	s.renderTemplate(w, "home.html", data)
}

// LoginHandler gère la connexion Spotify
func (s *Server) LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Rediriger vers la page d'authentification Spotify
	authURL := s.SpotifyAuth.GetAuthURL()
	http.Redirect(w, r, authURL, http.StatusSeeOther)
}

// CallbackHandler gère le callback Spotify après la connexion
func (s *Server) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	// Obtenir les paramètres de requête
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	// Valider l'état
	if state != s.SpotifyAuth.State {
		http.Error(w, "Paramètre d'état invalide", http.StatusBadRequest)
		return
	}

	// Échanger le code contre un jeton
	err := s.SpotifyAuth.Exchange(code)
	if err != nil {
		http.Error(w, "Échec lors de l'échange du code contre un jeton: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Rediriger vers la page d'accueil
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// LogoutHandler gère la déconnexion
func (s *Server) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Effacer les jetons
	s.SpotifyAuth.AccessToken = ""
	s.SpotifyAuth.RefreshToken = ""
	s.SpotifyAuth.Expiry = time.Time{}

	// Rediriger vers la page d'accueil
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) SearchHandler(w http.ResponseWriter, r *http.Request) {
	// Vérifier si l'utilisateur est connecté
	if !s.SpotifyAuth.IsTokenValid() {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Obtenir les paramètres de requête
	query := r.URL.Query().Get("q")
	typeParams := r.URL.Query()["type"] // Obtenir toutes les valeurs de type sous forme de slice
	limitStr := r.URL.Query().Get("limit")
	pageStr := r.URL.Query().Get("page")

	// Valeurs par défaut
	types := []string{"artist", "album", "track"}
	limit := 20
	page := 1

	// Analyser les paramètres de type
	if len(typeParams) > 0 {
		types = typeParams
		log.Printf("Types de recherche: %v", types)
	}

	// S'assurer qu'au moins un type est sélectionné
	if len(types) == 0 {
		types = []string{"artist", "album", "track"}
	}

	// Analyser le paramètre de limite
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// Analyser le paramètre de page
	if pageStr != "" {
		parsedPage, err := strconv.Atoi(pageStr)
		if err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	// Calculer le décalage
	offset := (page - 1) * limit

	// Si aucune requête, rendre la page de recherche vide
	if query == "" {
		data := PageData{
			Title:       "Recherche - MelodyExplorer",
			IsLoggedIn:  true,
			CurrentPage: "search",
			Query:       "",
		}
		s.renderTemplate(w, "search.html", data)
		return
	}

	// Initialiser des SearchResults vides
	searchResults := &spotify.SearchResults{}
	var searchErrors []error

	// Journaliser la recherche que nous sommes sur le point d'effectuer
	log.Printf("Exécution de la recherche: query=%s, types=%v, limit=%d, offset=%d", query, types, limit, offset)

	// Effectuer une recherche pour tous les types demandés en une seule fois
	if len(types) > 0 {
		var err error
		searchResults, err = s.SpotifyClient.Search(query, types, limit, offset)
		if err != nil {
			searchErrors = append(searchErrors, fmt.Errorf("erreur de recherche: %w", err))
		}
	}

	// Vérifier si la recherche a échoué
	if len(searchErrors) > 0 && searchResults == nil {
		errMessage := "Erreur lors de la recherche Spotify: "
		for i, err := range searchErrors {
			if i > 0 {
				errMessage += ", "
			}
			errMessage += err.Error()
		}

		data := PageData{
			Title:       "Recherche - MelodyExplorer",
			IsLoggedIn:  true,
			CurrentPage: "search",
			Query:       query,
			Error:       errMessage,
		}
		s.renderTemplate(w, "search.html", data)
		return
	}

	// Si nous n'avons pas obtenu de résultats, initialiser une structure vide
	if searchResults == nil {
		searchResults = &spotify.SearchResults{}
	}

	// Journaliser ce que nous avons reçu
	log.Printf("Résultats de recherche reçus:")
	if searchResults.Artists != nil {
		log.Printf("- Artistes: %d éléments", len(searchResults.Artists.Items))
	}
	if searchResults.Albums != nil {
		log.Printf("- Albums: %d éléments", len(searchResults.Albums.Items))
	}
	if searchResults.Tracks != nil {
		log.Printf("- Pistes: %d éléments", len(searchResults.Tracks.Items))
	}

	responseBody, _ := json.MarshalIndent(searchResults, "", "  ")
	log.Printf("Résultats bruts de la recherche:\n%s", string(responseBody))

	// S'assurer que les sections existent dans les résultats
	if searchResults.Artists == nil || searchResults.Artists.Items == nil {
		log.Printf("La section Artistes est nulle, création d'une section vide")
		searchResults.Artists = &spotify.ArtistResults{
			Items: []spotify.Artist{},
		}
	}
	if searchResults.Albums == nil || searchResults.Albums.Items == nil {
		log.Printf("La section Albums est nulle, création d'une section vide")
		searchResults.Albums = &spotify.AlbumResults{
			Items: []spotify.Album{},
		}
	}
	if searchResults.Tracks == nil || searchResults.Tracks.Items == nil {
		log.Printf("La section Pistes est nulle, création d'une section vide")
		searchResults.Tracks = &spotify.TrackResults{
			Items: []spotify.Track{},
		}
	}

	// Calculer la pagination pour les artistes
	var artistsPagination *PaginationData
	if searchResults.Artists != nil && len(searchResults.Artists.Items) > 0 {
		artistsPagination = &PaginationData{
			CurrentPage: page,
			TotalPages:  (searchResults.Artists.Total + limit - 1) / limit,
			TotalItems:  searchResults.Artists.Total,
			Limit:       limit,
			HasPrev:     page > 1,
			HasNext:     page*limit < searchResults.Artists.Total,
			PrevPage:    page - 1,
			NextPage:    page + 1,
		}
	}

	// Calculer la pagination pour les albums
	var albumsPagination *PaginationData
	if searchResults.Albums != nil && len(searchResults.Albums.Items) > 0 {
		albumsPagination = &PaginationData{
			CurrentPage: page,
			TotalPages:  (searchResults.Albums.Total + limit - 1) / limit,
			TotalItems:  searchResults.Albums.Total,
			Limit:       limit,
			HasPrev:     page > 1,
			HasNext:     page*limit < searchResults.Albums.Total,
			PrevPage:    page - 1,
			NextPage:    page + 1,
		}
	}

	// Calculer la pagination pour les pistes
	var tracksPagination *PaginationData
	if searchResults.Tracks != nil && len(searchResults.Tracks.Items) > 0 {
		tracksPagination = &PaginationData{
			CurrentPage: page,
			TotalPages:  (searchResults.Tracks.Total + limit - 1) / limit,
			TotalItems:  searchResults.Tracks.Total,
			Limit:       limit,
			HasPrev:     page > 1,
			HasNext:     page*limit < searchResults.Tracks.Total,
			PrevPage:    page - 1,
			NextPage:    page + 1,
		}
	}

	// Marquer les favoris
	favorites := s.FavoritesStorage.GetAll()
	favoriteMap := make(map[string]bool)
	for _, fav := range favorites {
		key := string(fav.Type) + ":" + fav.ID
		favoriteMap[key] = true
	}

	// Préparer les données de la page
	data := PageData{
		Title:       "Résultats de recherche - MelodyExplorer",
		IsLoggedIn:  true,
		CurrentPage: "search",
		Query:       query,
		Data: map[string]interface{}{
			"Results":           searchResults,
			"Types":             types,
			"Favorites":         favoriteMap,
			"ArtistsPagination": artistsPagination,
			"AlbumsPagination":  albumsPagination,
			"TracksPagination":  tracksPagination,
		},
	}

	// Rendre le template
	s.renderTemplate(w, "search.html", data)
}

// CollectionHandler gère la page de collection
func (s *Server) CollectionHandler(w http.ResponseWriter, r *http.Request) {
	// Vérifier si l'utilisateur est connecté
	if !s.SpotifyAuth.IsTokenValid() {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Obtenir les paramètres de requête
	typeParam := r.URL.Query().Get("type")
	genreParam := r.URL.Query().Get("genre")
	popularityParam := r.URL.Query().Get("popularity")
	yearParam := r.URL.Query().Get("year")
	limitStr := r.URL.Query().Get("limit")
	pageStr := r.URL.Query().Get("page")

	// Valeurs par défaut
	limit := 20
	page := 1

	// Analyser le paramètre de limite
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// Analyser le paramètre de page
	if pageStr != "" {
		parsedPage, err := strconv.Atoi(pageStr)
		if err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	// Calculer le décalage
	offset := (page - 1) * limit

	// Obtenir les genres disponibles pour le menu déroulant de filtrage
	genres, err := s.SpotifyClient.GetGenres()
	if err != nil {
		log.Printf("Erreur lors de l'obtention des genres: %v", err)
		genres = []string{}
	}

	// Créer une carte de filtres pour le template
	filters := map[string]string{
		"type":       typeParam,
		"genre":      genreParam,
		"popularity": popularityParam,
		"year":       yearParam,
	}

	// Préparer les paramètres de recherche pour la collection filtrée
	searchParams := "?"

	if typeParam != "" {
		searchParams += "type=" + typeParam + "&"
	}

	// Construire les chaînes de requête pour différents types de contenu
	var query string
	var artistQuery string

	// Personnaliser les requêtes en fonction des filtres
	if genreParam != "" {
		query += "genre:" + genreParam + " "
		artistQuery = "genre:" + genreParam + " "
		searchParams += "genre=" + genreParam + "&"
	}

	if yearParam != "" {
		query += "year:" + yearParam + " "
		searchParams += "year=" + yearParam + "&"
	}

	if popularityParam != "" {
		searchParams += "popularity=" + popularityParam + "&"
	}

	if query == "" {
		query = "year:2020-2024" // Par défaut pour les albums/pistes
	}

	if artistQuery == "" {
		artistQuery = "pop rock" // Par défaut pour les artistes - utiliser des genres populaires
	}

	searchParams += "limit=" + strconv.Itoa(limit) + "&page=" + strconv.Itoa(page)

	// Initialiser la structure des résultats
	var artistResults *spotify.SearchResults
	var albumResults *spotify.SearchResults
	var trackResults *spotify.SearchResults

	// Suivre si nous avons des résultats réussis
	var hasResults bool

	// Récupérer des données en fonction du filtre de type
	if typeParam == "" || typeParam == "artist" {
		// Obtenir des artistes
		log.Printf("Recherche d'artistes avec la requête: %s", artistQuery)
		artistResults, err = s.SpotifyClient.Search(artistQuery, []string{"artist"}, limit, offset)
		if err != nil {
			log.Printf("Erreur lors de la recherche d'artistes: %v", err)
		} else if artistResults != nil && artistResults.Artists != nil {
			hasResults = true
			log.Printf("Trouvé %d artistes", len(artistResults.Artists.Items))
		}
	}

	if typeParam == "" || typeParam == "album" {
		// Obtenir des albums
		albumResults, err = s.SpotifyClient.Search(query, []string{"album"}, limit, offset)
		if err != nil {
			log.Printf("Erreur lors de la recherche d'albums: %v", err)
		} else if albumResults != nil && albumResults.Albums != nil {
			hasResults = true
			log.Printf("Trouvé %d albums", len(albumResults.Albums.Items))
		}
	}

	if typeParam == "" || typeParam == "track" {
		// Obtenir des pistes
		trackResults, err = s.SpotifyClient.Search(query, []string{"track"}, limit, offset)
		if err != nil {
			log.Printf("Erreur lors de la recherche de pistes: %v", err)
		} else if trackResults != nil && trackResults.Tracks != nil {
			hasResults = true
			log.Printf("Trouvé %d pistes", len(trackResults.Tracks.Items))
		}
	}

	// Si aucun résultat trouvé, afficher une erreur
	if !hasResults {

		data := PageData{
			Title:       "Collection - MelodyExplorer",
			IsLoggedIn:  true,
			CurrentPage: "collection",
			Filters:     filters,
			Error:       "Aucun résultat ne correspond à vos filtres",
			Data: map[string]interface{}{
				"Genres": genres,
			},
		}
		s.renderTemplate(w, "collection.html", data)
		return
	}

	// Créer un objet de résultats combinés
	combinedResults := &spotify.SearchResults{}

	// Ajouter les résultats d'artistes
	if artistResults != nil && artistResults.Artists != nil {
		combinedResults.Artists = artistResults.Artists
	}

	// Ajouter les résultats d'albums
	if albumResults != nil && albumResults.Albums != nil {
		combinedResults.Albums = albumResults.Albums
	}

	// Ajouter les résultats de pistes
	if trackResults != nil && trackResults.Tracks != nil {
		combinedResults.Tracks = trackResults.Tracks
	}

	// Pour CollectionHandler:

	// Calculer le nombre total d'éléments (utiliser le nombre le plus élevé pour être cohérent)
	totalItems := 0
	totalArtists := 0
	totalAlbums := 0
	totalTracks := 0

	if combinedResults.Artists != nil {
		totalArtists = combinedResults.Artists.Total
	}
	if combinedResults.Albums != nil {
		totalAlbums = combinedResults.Albums.Total
	}
	if combinedResults.Tracks != nil {
		totalTracks = combinedResults.Tracks.Total
	}

	// Pour la cohérence, utiliser la valeur de comptage la plus élevée pour déterminer la pagination
	// De cette façon, même si un seul type de résultat est affiché, nous nous assurons que tous les éléments peuvent être accessibles
	if totalArtists > totalItems {
		totalItems = totalArtists
	}
	if totalAlbums > totalItems {
		totalItems = totalAlbums
	}
	if totalTracks > totalItems {
		totalItems = totalTracks
	}

	// Calculer le nombre total de pages - éviter de changer lorsque l'utilisateur navigue
	totalPages := (totalItems + limit - 1) / limit
	if totalPages < 1 {
		totalPages = 1
	}

	// S'assurer que nous n'affichons pas des numéros de page ridicules
	if totalPages > 100 {
		totalPages = 100 // Limite à un maximum de 100 pages
	}

	log.Printf("Page %d sur %d - Éléments: Artistes=%d, Albums=%d, Pistes=%d, Total le plus élevé=%d",
		page, totalPages, totalArtists, totalAlbums, totalTracks, totalItems)

	pagination := &PaginationData{
		CurrentPage: page,
		TotalPages:  totalPages,
		TotalItems:  totalItems,
		Limit:       limit,
		HasPrev:     page > 1,
		HasNext:     page < totalPages,
		PrevPage:    page - 1,
		NextPage:    page + 1,
	}

	// Ajouter un code similaire à SearchHandler aussi
	// Marquer les favoris
	favorites := s.FavoritesStorage.GetAll()
	favoriteMap := make(map[string]bool)
	for _, fav := range favorites {
		key := string(fav.Type) + ":" + fav.ID
		favoriteMap[key] = true
	}

	// Préparer les données de la page
	data := PageData{
		Title:       "Collection - MelodyExplorer",
		IsLoggedIn:  true,
		CurrentPage: "collection",
		Filters:     filters,
		Pagination:  pagination,
		Data: map[string]interface{}{
			"Results":      combinedResults,
			"Genres":       genres,
			"Favorites":    favoriteMap,
			"SearchParams": searchParams,
		},
	}

	// Rendre le template
	s.renderTemplate(w, "collection.html", data)
}

// ArtistHandler gère la page de détails de l'artiste
func (s *Server) ArtistHandler(w http.ResponseWriter, r *http.Request) {
	// Vérifier si l'utilisateur est connecté
	if !s.SpotifyAuth.IsTokenValid() {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Obtenir l'ID de l'artiste depuis l'URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Obtenir l'artiste depuis Spotify
	artist, err := s.SpotifyClient.GetArtist(id)
	if err != nil {
		s.ErrorHandler(w, r)
		log.Printf("Erreur lors de l'obtention de l'artiste: %v", err)
		return
	}

	// Obtenir les albums de l'artiste
	albums, err := s.SpotifyClient.GetArtistAlbums(id, 20, 0)
	if err != nil {
		s.ErrorHandler(w, r)
		log.Printf("Erreur lors de l'obtention des albums de l'artiste: %v", err)
		return
	}

	// Vérifier si l'artiste est dans les favoris
	isFavorite := s.FavoritesStorage.Contains(id, models.FavoriteTypeArtist)

	// Préparer les données de la page
	data := PageData{
		Title:       artist.Name + " - MelodyExplorer",
		IsLoggedIn:  true,
		CurrentPage: "artist",
		Pagination:  nil, // Définir explicitement à nil
		Data: map[string]interface{}{
			"Artist":     artist,
			"Albums":     albums,
			"IsFavorite": isFavorite,
		},
	}

	// Rendre le template
	s.renderTemplate(w, "details.html", data)
}

// AlbumHandler gère la page de détails de l'album
func (s *Server) AlbumHandler(w http.ResponseWriter, r *http.Request) {
	// Vérifier si l'utilisateur est connecté
	if !s.SpotifyAuth.IsTokenValid() {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Obtenir l'ID de l'album depuis l'URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Obtenir l'album depuis Spotify
	album, err := s.SpotifyClient.GetAlbum(id)
	if err != nil {
		s.ErrorHandler(w, r)
		log.Printf("Erreur lors de l'obtention de l'album: %v", err)
		return
	}

	// Obtenir les pistes de l'album
	tracks, err := s.SpotifyClient.GetAlbumTracks(id, 50, 0)
	if err != nil {
		s.ErrorHandler(w, r)
		log.Printf("Erreur lors de l'obtention des pistes de l'album: %v", err)
		return
	}

	// Vérifier si l'album est dans les favoris
	isFavorite := s.FavoritesStorage.Contains(id, models.FavoriteTypeAlbum)

	// Préparer les données de la page
	data := PageData{
		Title:       album.Name + " - MelodyExplorer",
		IsLoggedIn:  true,
		CurrentPage: "album",
		Pagination:  nil, // Définir explicitement à nil
		Data: map[string]interface{}{
			"Album":      album,
			"Tracks":     tracks,
			"IsFavorite": isFavorite,
		},
	}

	// Rendre le template
	s.renderTemplate(w, "details.html", data)
}

// TrackHandler gère la page de détails de la piste
func (s *Server) TrackHandler(w http.ResponseWriter, r *http.Request) {
	// Vérifier si l'utilisateur est connecté
	if !s.SpotifyAuth.IsTokenValid() {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Obtenir l'ID de la piste depuis l'URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Obtenir la piste depuis Spotify
	track, err := s.SpotifyClient.GetTrack(id)
	if err != nil {
		s.ErrorHandler(w, r)
		log.Printf("Erreur lors de l'obtention de la piste: %v", err)
		return
	}

	// Vérifier si la piste est dans les favoris
	isFavorite := s.FavoritesStorage.Contains(id, models.FavoriteTypeTrack)

	// Préparer les données de la page
	data := PageData{
		Title:       track.Name + " - MelodyExplorer",
		IsLoggedIn:  true,
		CurrentPage: "track",
		Pagination:  nil, // Définir explicitement à nil
		Data: map[string]interface{}{
			"Track":      track,
			"IsFavorite": isFavorite,
		},
	}

	// Rendre le template
	s.renderTemplate(w, "details.html", data)
}

func (s *Server) AddFavoriteHandler(w http.ResponseWriter, r *http.Request) {
	// Vérifier si l'utilisateur est connecté
	if !s.SpotifyAuth.IsTokenValid() {
		http.Error(w, "Non autorisé", http.StatusUnauthorized)
		return
	}

	// Analyser le corps de la requête
	var req struct {
		ID       string `json:"id"`
		Type     string `json:"type"`
		Name     string `json:"name"`
		ImageURL string `json:"image_url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Corps de requête invalide", http.StatusBadRequest)
		log.Printf("Erreur lors de l'analyse de la requête d'ajout de favori: %v", err)
		return
	}

	// Valider le type
	var itemType models.FavoriteType
	switch req.Type {
	case "artist":
		itemType = models.FavoriteTypeArtist
	case "album":
		itemType = models.FavoriteTypeAlbum
	case "track":
		itemType = models.FavoriteTypeTrack
	default:
		http.Error(w, "Type invalide", http.StatusBadRequest)
		log.Printf("Type de favori invalide: %s", req.Type)
		return
	}

	log.Printf("Ajout aux favoris: %s (%s - %s)", req.Name, req.Type, req.ID)

	// Créer un élément favori
	item := models.FavoriteItem{
		ID:       req.ID,
		Type:     itemType,
		Name:     req.Name,
		ImageURL: req.ImageURL,
		AddedAt:  time.Now(),
	}

	// Ajouter aux favoris et sauvegarder dans le fichier
	if err := s.FavoritesStorage.Add(item); err != nil {
		http.Error(w, "Échec lors de l'ajout aux favoris: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Erreur lors de l'ajout du favori: %v", err)
		return
	}

	// Retourner succès
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

// RemoveFavoriteHandler gère la suppression d'éléments des favoris
func (s *Server) RemoveFavoriteHandler(w http.ResponseWriter, r *http.Request) {
	// Vérifier si l'utilisateur est connecté
	if !s.SpotifyAuth.IsTokenValid() {
		http.Error(w, "Non autorisé", http.StatusUnauthorized)
		return
	}

	// Analyser le corps de la requête
	var req struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Corps de requête invalide", http.StatusBadRequest)
		log.Printf("Erreur lors de l'analyse de la requête de suppression de favori: %v", err)
		return
	}

	// Valider le type
	var itemType models.FavoriteType
	switch req.Type {
	case "artist":
		itemType = models.FavoriteTypeArtist
	case "album":
		itemType = models.FavoriteTypeAlbum
	case "track":
		itemType = models.FavoriteTypeTrack
	default:
		http.Error(w, "Type invalide", http.StatusBadRequest)
		log.Printf("Type de favori invalide pour la suppression: %s", req.Type)
		return
	}

	log.Printf("Suppression du favori: %s (%s)", req.ID, req.Type)

	// Supprimer des favoris et sauvegarder dans le fichier
	if err := s.FavoritesStorage.Remove(req.ID, itemType); err != nil {
		http.Error(w, "Échec lors de la suppression des favoris: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Erreur lors de la suppression du favori: %v", err)
		return
	}

	// Retourner succès
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

// FavoritesHandler gère la page des favoris
func (s *Server) FavoritesHandler(w http.ResponseWriter, r *http.Request) {
	// Vérifier si l'utilisateur est connecté
	if !s.SpotifyAuth.IsTokenValid() {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Obtenir les favoris
	favorites := s.FavoritesStorage.GetAll()

	// Organiser les favoris par type
	artists := s.FavoritesStorage.GetByType(models.FavoriteTypeArtist)
	albums := s.FavoritesStorage.GetByType(models.FavoriteTypeAlbum)
	tracks := s.FavoritesStorage.GetByType(models.FavoriteTypeTrack)

	// Préparer les données de la page
	data := PageData{
		Title:       "Mes Favoris - MelodyExplorer",
		IsLoggedIn:  true,
		CurrentPage: "favorites",
		Data: map[string]interface{}{
			"Favorites": favorites,
			"Artists":   artists,
			"Albums":    albums,
			"Tracks":    tracks,
		},
	}

	// Rendre le template
	s.renderTemplate(w, "favorites.html", data)
}

// CategoryHandler gère la page de catégorie/genre
func (s *Server) CategoryHandler(w http.ResponseWriter, r *http.Request) {
	// Vérifier si l'utilisateur est connecté
	if !s.SpotifyAuth.IsTokenValid() {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Obtenir le genre depuis l'URL
	vars := mux.Vars(r)
	genre := vars["genre"]

	// Obtenir les paramètres de pagination
	limitStr := r.URL.Query().Get("limit")
	pageStr := r.URL.Query().Get("page")

	// Valeurs par défaut
	limit := 20
	page := 1

	// Analyser le paramètre de limite
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// Analyser le paramètre de page
	if pageStr != "" {
		parsedPage, err := strconv.Atoi(pageStr)
		if err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	// Calculer le décalage
	offset := (page - 1) * limit

	// Rechercher le genre avec une meilleure requête pour obtenir plus de résultats
	query := "genre:" + genre
	// Ajouter quelques termes de filtre populaires pour obtenir de meilleurs résultats
	if page > 5 {
		// Pour les pages ultérieures, ajouter des plages d'années pour obtenir des résultats plus diversifiés
		yearStart := 2020 - ((page / 5) * 10)
		yearEnd := yearStart + 9
		if yearStart < 1950 {
			yearStart = 1950
		}
		query += " year:" + strconv.Itoa(yearStart) + "-" + strconv.Itoa(yearEnd)
	}

	log.Printf("Requête de recherche de catégorie: %s, page: %d, décalage: %d", query, page, offset)

	// Essayer avec une limite plus élevée pour s'assurer d'obtenir suffisamment de résultats
	searchLimit := limit * 2
	results, err := s.SpotifyClient.Search(query, []string{"track"}, searchLimit, offset)
	if err != nil {
		http.Error(w, "Erreur lors de la recherche Spotify: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Calculer la pagination avec des pages totales fixes pour assurer une navigation cohérente
	totalItems := 1000 // Définir une limite fixe raisonnable
	if results.Tracks != nil && results.Tracks.Total > 0 && results.Tracks.Total < totalItems {
		totalItems = results.Tracks.Total
	}

	totalPages := (totalItems + limit - 1) / limit

	// S'assurer que nous affichons au moins 20 pages si nous avons des résultats
	if totalPages < 20 && results.Tracks != nil && len(results.Tracks.Items) > 0 {
		totalPages = 20
	}

	pagination := &PaginationData{
		CurrentPage: page,
		TotalPages:  totalPages,
		TotalItems:  totalItems,
		Limit:       limit,
		HasPrev:     page > 1,
		HasNext:     page < totalPages,
		PrevPage:    page - 1,
		NextPage:    page + 1,
	}

	// Marquer les favoris
	favorites := s.FavoritesStorage.GetAll()
	favoriteMap := make(map[string]bool)
	for _, fav := range favorites {
		key := string(fav.Type) + ":" + fav.ID
		favoriteMap[key] = true
	}

	// Préparer les données de la page
	data := PageData{
		Title:       genre + " Musique - MelodyExplorer",
		IsLoggedIn:  true,
		CurrentPage: "category",
		Pagination:  pagination,
		Data: map[string]interface{}{
			"Genre":     genre,
			"Results":   results,
			"Favorites": favoriteMap,
		},
	}

	// Rendre le template
	s.renderTemplate(w, "category.html", data)
}

// AboutHandler gère la page à propos
func (s *Server) AboutHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title:       "À propos - MelodyExplorer",
		IsLoggedIn:  s.SpotifyAuth.IsTokenValid(),
		CurrentPage: "about",
	}

	s.renderTemplate(w, "about.html", data)
}

// ErrorHandler gère les erreurs
func (s *Server) ErrorHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title:       "Erreur - MelodyExplorer",
		IsLoggedIn:  s.SpotifyAuth.IsTokenValid(),
		CurrentPage: "error",
		Error:       "La page que vous avez demandée est introuvable",
	}

	s.renderTemplate(w, "error.html", data)
}
