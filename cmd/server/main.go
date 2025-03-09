package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"github.com/yourusername/melody-explorer/internal/api"
)

func main() {
	// Charger les variables d'environnement depuis le fichier .env
	if err := godotenv.Load(); err != nil {
		log.Println("Avertissement: fichier .env introuvable")
	}

	// Récupérer le port depuis la variable d'environnement ou utiliser la valeur par défaut
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Déterminer le répertoire racine du projet
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Erreur lors de l'obtention du répertoire de travail: %v", err)
	}

	// Afficher le répertoire de travail pour le débogage
	fmt.Printf("Répertoire de travail: %s\n", wd)

	// Résoudre les chemins
	templatesDir := filepath.Join(wd, "templates")
	staticDir := filepath.Join(wd, "static")
	dataDir := filepath.Join(wd, "data")

	// Afficher le répertoire de données pour le débogage
	fmt.Printf("Répertoire des données: %s\n", dataDir)

	// Afficher le chemin de favorites.json pour le débogage
	favoritesPath := filepath.Join(dataDir, "favorites.json")
	fmt.Printf("Chemin du fichier des favoris: %s\n", favoritesPath)

	// Vérifier si le fichier favorites.json existe
	if _, err := os.Stat(favoritesPath); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("favorites.json n'existe pas!")
		} else {
			fmt.Printf("Erreur lors de la vérification de favorites.json: %v\n", err)
		}
	} else {
		fmt.Println("favorites.json existe")

		// Essayer de lire le fichier et afficher son contenu
		data, err := os.ReadFile(favoritesPath)
		if err != nil {
			fmt.Printf("Erreur lors de la lecture de favorites.json: %v\n", err)
		} else {
			fmt.Printf("Contenu de favorites.json: %q\n", string(data))
		}
	}

	// Créer le répertoire de données s'il n'existe pas
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		fmt.Println("Création du répertoire de données...")
		if err := os.MkdirAll(dataDir, 0755); err != nil {
			log.Fatalf("Erreur lors de la création du répertoire de données: %v", err)
		}
	}

	// Créer un fichier favorites.json propre avec un encodage approprié
	fmt.Println("Création d'un fichier favorites.json propre...")
	emptyArray := []byte("[]")
	if err := os.WriteFile(favoritesPath, emptyArray, 0644); err != nil {
		log.Fatalf("Erreur lors de la création de favorites.json: %v", err)
	}

	fmt.Println("Création réussie du fichier favorites.json propre")

	// Créer le serveur
	server, err := api.NewServer(templatesDir, staticDir, dataDir)
	if err != nil {
		log.Fatalf("Erreur lors de la création du serveur: %v", err)
	}

	// Créer le serveur HTTP
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      server.Router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Démarrer le serveur dans une goroutine
	go func() {
		log.Printf("Serveur à l'écoute sur le port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erreur lors du démarrage du serveur: %v", err)
		}
	}()

	// Maintenir la goroutine principale active
	select {}
}
