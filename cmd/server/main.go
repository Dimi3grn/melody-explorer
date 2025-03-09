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
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Determine project root directory
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting working directory: %v", err)
	}

	// Print the working directory for debugging
	fmt.Printf("Working directory: %s\n", wd)

	// Resolve paths
	templatesDir := filepath.Join(wd, "templates")
	staticDir := filepath.Join(wd, "static")
	dataDir := filepath.Join(wd, "data")

	// Print the data directory for debugging
	fmt.Printf("Data directory: %s\n", dataDir)

	// Print the favorites.json path for debugging
	favoritesPath := filepath.Join(dataDir, "favorites.json")
	fmt.Printf("Favorites file path: %s\n", favoritesPath)

	// Check if the favorites.json file exists
	if _, err := os.Stat(favoritesPath); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("favorites.json does not exist!")
		} else {
			fmt.Printf("Error checking favorites.json: %v\n", err)
		}
	} else {
		fmt.Println("favorites.json exists")

		// Try to read the file and print its contents
		data, err := os.ReadFile(favoritesPath)
		if err != nil {
			fmt.Printf("Error reading favorites.json: %v\n", err)
		} else {
			fmt.Printf("Content of favorites.json: %q\n", string(data))
		}
	}

	// Create data directory if it doesn't exist
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		fmt.Println("Creating data directory...")
		if err := os.MkdirAll(dataDir, 0755); err != nil {
			log.Fatalf("Error creating data directory: %v", err)
		}
	}

	// Create a clean favorites.json file with proper encoding
	fmt.Println("Creating a clean favorites.json...")
	emptyArray := []byte("[]")
	if err := os.WriteFile(favoritesPath, emptyArray, 0644); err != nil {
		log.Fatalf("Error creating favorites.json: %v", err)
	}

	fmt.Println("Created clean favorites.json successfully")

	// Create server
	server, err := api.NewServer(templatesDir, staticDir, dataDir)
	if err != nil {
		log.Fatalf("Error creating server: %v", err)
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      server.Router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server listening on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Keep the main goroutine alive
	select {}
}
