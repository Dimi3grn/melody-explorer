package api

import (
	"net/http"
)

// initializeRoutes configure toutes les routes pour le serveur
func (s *Server) initializeRoutes() {
	// Appliquer le middleware d'authentification
	s.Router.Use(s.SpotifyAuth.AuthMiddleware)

	// Fichiers statiques
	fs := http.FileServer(http.Dir(s.StaticDir))
	s.Router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	// Routes d'authentification
	s.Router.HandleFunc("/login", s.LoginHandler).Methods("GET")
	s.Router.HandleFunc("/callback", s.CallbackHandler).Methods("GET")
	s.Router.HandleFunc("/logout", s.LogoutHandler).Methods("GET")

	// Routes des pages
	s.Router.HandleFunc("/", s.HomeHandler).Methods("GET")
	s.Router.HandleFunc("/search", s.SearchHandler).Methods("GET")
	s.Router.HandleFunc("/collection", s.CollectionHandler).Methods("GET")
	s.Router.HandleFunc("/artist/{id}", s.ArtistHandler).Methods("GET")
	s.Router.HandleFunc("/album/{id}", s.AlbumHandler).Methods("GET")
	s.Router.HandleFunc("/track/{id}", s.TrackHandler).Methods("GET")
	s.Router.HandleFunc("/favorites", s.FavoritesHandler).Methods("GET")
	s.Router.HandleFunc("/category/{genre}", s.CategoryHandler).Methods("GET")
	s.Router.HandleFunc("/about", s.AboutHandler).Methods("GET")
	s.Router.HandleFunc("/recommandation", s.RecommendationHandler).Methods("GET")

	// Routes API
	s.Router.HandleFunc("/api/favorites/add", s.AddFavoriteHandler).Methods("POST")
	s.Router.HandleFunc("/api/favorites/remove", s.RemoveFavoriteHandler).Methods("POST")

	// Gestionnaire d'erreur (404)
	s.Router.NotFoundHandler = http.HandlerFunc(s.ErrorHandler)
}
