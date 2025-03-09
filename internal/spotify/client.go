package spotify

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	baseURL = "https://api.spotify.com/v1"
)

// Client est un client d'API Spotify
type Client struct {
	Auth *Auth
}

// NewClient crée un nouveau client d'API Spotify
func NewClient(auth *Auth) *Client {
	return &Client{
		Auth: auth,
	}
}

// makeRequest effectue une requête HTTP vers l'API Spotify
func (c *Client) makeRequest(method, endpoint string, params url.Values) ([]byte, error) {
	// S'assurer que nous avons un jeton valide
	if err := c.Auth.EnsureValidToken(); err != nil {
		return nil, err
	}

	// Construire l'URL
	apiURL := baseURL + endpoint
	if params != nil && len(params) > 0 {
		apiURL += "?" + params.Encode()
	}

	// Journaliser l'URL de la requête (pour le débogage)
	log.Printf("Requête vers : %s", apiURL)

	// Créer la requête
	req, err := http.NewRequest(method, apiURL, nil)
	if err != nil {
		log.Printf("Erreur lors de la création de la requête : %v", err)
		return nil, err
	}

	// Ajouter l'en-tête d'autorisation
	req.Header.Add("Authorization", "Bearer "+c.Auth.AccessToken)

	// Journaliser le jeton (10 premiers caractères pour la sécurité)
	tokenPreview := c.Auth.AccessToken
	if len(tokenPreview) > 10 {
		tokenPreview = tokenPreview[:10] + "..."
	}
	log.Printf("Utilisation du jeton : %s", tokenPreview)

	// Effectuer la requête
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Erreur lors de l'exécution de la requête : %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Journaliser le code d'état de la réponse
	log.Printf("Statut de la réponse : %d %s", resp.StatusCode, resp.Status)

	// Vérifier les erreurs
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		errorMsg := fmt.Sprintf("erreur API Spotify (%d) : %s", resp.StatusCode, string(body))
		log.Printf("Réponse d'erreur API : %s", string(body))
		return nil, fmt.Errorf(errorMsg)
	}

	// Lire le corps de la réponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Erreur lors de la lecture du corps de la réponse : %v", err)
		return nil, err
	}

	// Journaliser les 100 premiers caractères de la réponse (pour le débogage)
	responsePreview := string(body)
	if len(responsePreview) > 100 {
		responsePreview = responsePreview[:100] + "..."
	}
	log.Printf("Aperçu du corps de la réponse : %s", responsePreview)

	return body, nil
}

// SearchResults contient les résultats de recherche de Spotify
type SearchResults struct {
	Artists *ArtistResults `json:"artists,omitempty"`
	Albums  *AlbumResults  `json:"albums,omitempty"`
	Tracks  *TrackResults  `json:"tracks,omitempty"`
}

// ArtistResults contient les résultats de recherche d'artistes
type ArtistResults struct {
	Items  []Artist `json:"items"`
	Total  int      `json:"total"`
	Limit  int      `json:"limit"`
	Offset int      `json:"offset"`
}

// AlbumResults contient les résultats de recherche d'albums
type AlbumResults struct {
	Items  []Album `json:"items"`
	Total  int     `json:"total"`
	Limit  int     `json:"limit"`
	Offset int     `json:"offset"`
}

// TrackResults contient les résultats de recherche de pistes
type TrackResults struct {
	Items  []Track `json:"items"`
	Total  int     `json:"total"`
	Limit  int     `json:"limit"`
	Offset int     `json:"offset"`
}

// Artist représente un artiste Spotify
type Artist struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Popularity   int               `json:"popularity"`
	Genres       []string          `json:"genres"`
	Images       []Image           `json:"images"`
	Followers    Followers         `json:"followers"`
	ExternalURLs map[string]string `json:"external_urls"`
}

// Album représente un album Spotify
type Album struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	ReleaseDate  string            `json:"release_date"`
	Images       []Image           `json:"images"`
	Artists      []Artist          `json:"artists"`
	TotalTracks  int               `json:"total_tracks"`
	AlbumType    string            `json:"album_type"`
	ExternalURLs map[string]string `json:"external_urls"`
}

// Track représente une piste Spotify
type Track struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Album        Album             `json:"album"`
	Artists      []Artist          `json:"artists"`
	Duration     int               `json:"duration_ms"`
	Popularity   int               `json:"popularity"`
	ExternalURLs map[string]string `json:"external_urls"`
}

// Image représente une image Spotify
type Image struct {
	URL    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

// Followers représente les informations sur les abonnés Spotify
type Followers struct {
	Total int `json:"total"`
}

// Search recherche des artistes, des albums et des pistes
func (c *Client) Search(query string, types []string, limit, offset int) (*SearchResults, error) {
	params := url.Values{}
	params.Add("q", query)
	params.Add("type", strings.Join(types, ","))
	params.Add("limit", strconv.Itoa(limit))
	params.Add("offset", strconv.Itoa(offset))

	body, err := c.makeRequest("GET", "/search", params)
	if err != nil {
		return nil, err
	}

	var results SearchResults
	if err := json.Unmarshal(body, &results); err != nil {
		return nil, err
	}

	return &results, nil
}

// GetArtist récupère un seul artiste par ID
func (c *Client) GetArtist(id string) (*Artist, error) {
	body, err := c.makeRequest("GET", "/artists/"+id, nil)
	if err != nil {
		return nil, err
	}

	var artist Artist
	if err := json.Unmarshal(body, &artist); err != nil {
		return nil, err
	}

	return &artist, nil
}

// GetArtistAlbums récupère les albums d'un artiste
func (c *Client) GetArtistAlbums(id string, limit, offset int) (*AlbumResults, error) {
	params := url.Values{}
	params.Add("limit", strconv.Itoa(limit))
	params.Add("offset", strconv.Itoa(offset))

	body, err := c.makeRequest("GET", "/artists/"+id+"/albums", params)
	if err != nil {
		return nil, err
	}

	var results AlbumResults
	if err := json.Unmarshal(body, &results); err != nil {
		return nil, err
	}

	return &results, nil
}

// GetAlbum récupère un seul album par ID
func (c *Client) GetAlbum(id string) (*Album, error) {
	body, err := c.makeRequest("GET", "/albums/"+id, nil)
	if err != nil {
		return nil, err
	}

	var album Album
	if err := json.Unmarshal(body, &album); err != nil {
		return nil, err
	}

	return &album, nil
}

// GetAlbumTracks récupère les pistes d'un album
func (c *Client) GetAlbumTracks(id string, limit, offset int) (*TrackResults, error) {
	params := url.Values{}
	params.Add("limit", strconv.Itoa(limit))
	params.Add("offset", strconv.Itoa(offset))

	body, err := c.makeRequest("GET", "/albums/"+id+"/tracks", params)
	if err != nil {
		return nil, err
	}

	var results TrackResults
	if err := json.Unmarshal(body, &results); err != nil {
		return nil, err
	}

	return &results, nil
}

// GetTrack récupère une seule piste par ID
func (c *Client) GetTrack(id string) (*Track, error) {
	body, err := c.makeRequest("GET", "/tracks/"+id, nil)
	if err != nil {
		return nil, err
	}

	var track Track
	if err := json.Unmarshal(body, &track); err != nil {
		return nil, err
	}

	return &track, nil
}

// GetGenres récupère les catégories musicales de Spotify
// Comme le point de terminaison genre-seeds peut être obsolète, nous utilisons les catégories à la place
func (c *Client) GetGenres() ([]string, error) {
	log.Println("Récupération des catégories musicales comme genres")

	// Utiliser le point de terminaison browse/categories
	params := url.Values{}
	params.Add("limit", "50") // Récupérer jusqu'à 50 catégories

	body, err := c.makeRequest("GET", "/browse/categories", params)
	if err != nil {
		log.Printf("Erreur lors de la récupération des catégories : %v", err)

		// Si l'appel API échoue, renvoyer une liste fixe de genres populaires
		return []string{
			"pop", "rock", "hip-hop", "jazz", "electronic",
			"classical", "country", "reggae", "blues",
			"metal", "indie", "folk", "r&b", "latin",
			"alternative", "dance", "ambient", "punk",
		}, nil
	}

	// Analyser la réponse des catégories
	var response struct {
		Categories struct {
			Items []struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"items"`
		} `json:"categories"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		log.Printf("Erreur lors de l'analyse de la réponse des catégories : %v", err)
		return nil, err
	}

	// Extraire les noms des catégories comme genres
	genres := make([]string, 0, len(response.Categories.Items))
	for _, item := range response.Categories.Items {
		genres = append(genres, item.Name)
	}

	return genres, nil
}
func (c *Client) GetRecommendations(seedGenres []string, limit int) ([]Track, error) {
	log.Println("Récupération des recommandations musicales en utilisant le point de terminaison new-releases")

	// Nous utiliserons les nouvelles sorties comme source de pistes recommandées
	params := url.Values{}
	params.Add("limit", strconv.Itoa(limit))

	// Obtenir les nouvelles sorties
	body, err := c.makeRequest("GET", "/browse/new-releases", params)
	if err != nil {
		log.Printf("Erreur lors de la récupération des nouvelles sorties : %v", err)
		return nil, err
	}

	// Analyser la réponse des nouvelles sorties
	var albumsResponse struct {
		Albums struct {
			Items []Album `json:"items"`
		} `json:"albums"`
	}

	if err := json.Unmarshal(body, &albumsResponse); err != nil {
		log.Printf("Erreur lors de l'analyse de la réponse des nouvelles sorties : %v", err)
		return nil, err
	}

	// Nous devons obtenir les pistes de ces albums
	tracks := make([]Track, 0, limit)
	for _, album := range albumsResponse.Albums.Items {
		// Obtenir les pistes pour cet album
		albumTracksParams := url.Values{}
		albumTracksParams.Add("limit", "1") // Obtenir seulement 1 piste par album

		albumTracksBody, err := c.makeRequest("GET", "/albums/"+album.ID+"/tracks", albumTracksParams)
		if err != nil {
			log.Printf("Erreur lors de la récupération des pistes pour l'album %s : %v", album.ID, err)
			continue
		}

		var tracksResponse struct {
			Items []Track `json:"items"`
		}

		if err := json.Unmarshal(albumTracksBody, &tracksResponse); err != nil {
			log.Printf("Erreur lors de l'analyse de la réponse des pistes d'album : %v", err)
			continue
		}

		// Ajouter les pistes de cet album
		for _, track := range tracksResponse.Items {
			// Remplir les informations de l'album qui pourraient manquer
			track.Album = album
			tracks = append(tracks, track)

			// Arrêter si nous avons atteint notre limite
			if len(tracks) >= limit {
				break
			}
		}

		// Arrêter si nous avons atteint notre limite
		if len(tracks) >= limit {
			break
		}
	}

	return tracks, nil
}
