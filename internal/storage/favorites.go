package storage

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/yourusername/melody-explorer/internal/models"
)

// FavoritesStorage gère le stockage des favoris
type FavoritesStorage struct {
	filename  string
	favorites *models.Favorites
	mu        sync.Mutex
}

// NewFavoritesStorage crée un nouveau FavoritesStorage
func NewFavoritesStorage(dataDir string) (*FavoritesStorage, error) {
	// S'assurer que le répertoire de données existe
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, err
	}

	filename := filepath.Join(dataDir, "favorites.json")
	favorites := models.NewFavorites()

	// Créer un fichier de favoris vide s'il n'existe pas
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.Printf("Création d'un nouveau fichier de favoris à %s", filename)
		emptyJSON := []byte("[]")
		if err := os.WriteFile(filename, emptyJSON, 0644); err != nil {
			return nil, err
		}
	}

	storage := &FavoritesStorage{
		filename:  filename,
		favorites: favorites,
	}

	// Charger les favoris existants
	if err := storage.Load(); err != nil {
		log.Printf("Erreur lors du chargement des favoris, démarrage avec une liste vide : %v", err)
	} else {
		log.Printf("Chargement de %d favoris depuis %s", len(storage.favorites.Items), filename)
	}

	return storage, nil
}

// Add ajoute un élément favori
func (s *FavoritesStorage) Add(item models.FavoriteItem) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.favorites.Add(item)

	// Sauvegarder immédiatement
	err := s.save()
	if err != nil {
		log.Printf("Erreur lors de la sauvegarde des favoris : %v", err)
	} else {
		log.Printf("Favori sauvegardé : %s (%s)", item.Name, item.Type)
	}

	return err
}

// Remove supprime un élément favori
func (s *FavoritesStorage) Remove(id string, itemType models.FavoriteType) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.favorites.Remove(id, itemType)

	// Sauvegarder immédiatement
	err := s.save()
	if err != nil {
		log.Printf("Erreur lors de la sauvegarde des favoris après suppression : %v", err)
	} else {
		log.Printf("Favori supprimé et sauvegardé : %s (%s)", id, itemType)
	}

	return err
}

// GetAll renvoie tous les éléments favoris
func (s *FavoritesStorage) GetAll() []models.FavoriteItem {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.favorites.Get()
}

// GetByType renvoie les éléments favoris d'un type spécifique
func (s *FavoritesStorage) GetByType(itemType models.FavoriteType) []models.FavoriteItem {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.favorites.GetByType(itemType)
}

// Contains vérifie si un élément est dans les favoris
func (s *FavoritesStorage) Contains(id string, itemType models.FavoriteType) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.favorites.Contains(id, itemType)
}

// save sauvegarde les favoris dans le fichier
func (s *FavoritesStorage) save() error {
	log.Printf("Sauvegarde des favoris dans %s", s.filename)
	data, err := json.MarshalIndent(s.favorites.Items, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filename, data, 0644)
}

// Load charge les favoris depuis le fichier
func (s *FavoritesStorage) Load() error {
	log.Printf("Chargement des favoris depuis %s", s.filename)
	data, err := os.ReadFile(s.filename)
	if err != nil {
		// Si le fichier n'existe pas, initialiser avec un tableau vide
		if os.IsNotExist(err) {
			s.favorites.Items = []models.FavoriteItem{}
			return nil
		}
		return err
	}

	// Si le fichier est vide ou ne contient que des espaces, initialiser avec un tableau vide
	if len(data) == 0 || string(data) == "[]" {
		s.favorites.Items = []models.FavoriteItem{}
		return nil
	}

	var items []models.FavoriteItem
	if err := json.Unmarshal(data, &items); err != nil {
		return err
	}

	s.favorites.Items = items
	return nil
}
