package api

import (
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware enregistre des informations sur chaque requête
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Appeler le gestionnaire suivant
		next.ServeHTTP(w, r)

		// Journaliser la requête
		log.Printf(
			"%s %s %s %s",
			r.RemoteAddr,
			r.Method,
			r.URL.Path,
			time.Since(start),
		)
	})
}

// RecoverMiddleware récupère les paniques et journalise l'erreur
func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panique: %v", err)
				http.Error(w, "Erreur Interne du Serveur", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// CORSMiddleware ajoute des en-têtes CORS aux réponses
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Définir les en-têtes CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")

		// Gérer les requêtes préliminaires
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Appeler le gestionnaire suivant
		next.ServeHTTP(w, r)
	})
}

// CacheControlMiddleware ajoute des en-têtes de contrôle de cache pour les ressources statiques
func CacheControlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ajouter des en-têtes de contrôle de cache pour les ressources statiques
		if r.URL.Path[:8] == "/static/" {
			w.Header().Set("Cache-Control", "public, max-age=86400") // 1 jour
		}

		// Appeler le gestionnaire suivant
		next.ServeHTTP(w, r)
	})
}
