package models

import (
	"sync"
	"time"
)

// FavoriteType définit le type d'élément favori
type FavoriteType string

const (
	// FavoriteTypeArtist représente un artiste favori
	FavoriteTypeArtist FavoriteType = "artist"
	// FavoriteTypeAlbum représente un album favori
	FavoriteTypeAlbum FavoriteType = "album"
	// FavoriteTypeTrack représente une piste favorite
	FavoriteTypeTrack FavoriteType = "track"
)

// FavoriteItem représente un élément favori
type FavoriteItem struct {
	ID       string       `json:"id"`
	Type     FavoriteType `json:"type"`
	Name     string       `json:"name"`
	ImageURL string       `json:"image_url"`
	AddedAt  time.Time    `json:"added_at"`
}

// Favorites représente une collection d'éléments favoris
type Favorites struct {
	Items []FavoriteItem `json:"items"`
	mu    sync.Mutex
}

// NewFavorites crée une nouvelle instance de Favorites
func NewFavorites() *Favorites {
	return &Favorites{
		Items: []FavoriteItem{},
	}
}

// Add ajoute un élément favori
func (f *Favorites) Add(item FavoriteItem) {
	f.mu.Lock()
	defer f.mu.Unlock()

	// Vérifier si l'élément existe déjà
	for i, existing := range f.Items {
		if existing.ID == item.ID && existing.Type == item.Type {
			// Mettre à jour l'élément
			f.Items[i] = item
			return
		}
	}

	// Ajouter l'élément
	if item.AddedAt.IsZero() {
		item.AddedAt = time.Now()
	}
	f.Items = append(f.Items, item)
}

// Remove supprime un élément favori
func (f *Favorites) Remove(id string, itemType FavoriteType) {
	f.mu.Lock()
	defer f.mu.Unlock()

	var newItems []FavoriteItem
	for _, item := range f.Items {
		if item.ID != id || item.Type != itemType {
			newItems = append(newItems, item)
		}
	}
	f.Items = newItems
}

// Get renvoie tous les éléments favoris
func (f *Favorites) Get() []FavoriteItem {
	f.mu.Lock()
	defer f.mu.Unlock()

	// Renvoyer une copie pour éviter les conditions de course
	items := make([]FavoriteItem, len(f.Items))
	copy(items, f.Items)
	return items
}

// GetByType renvoie les éléments favoris d'un type spécifique
func (f *Favorites) GetByType(itemType FavoriteType) []FavoriteItem {
	f.mu.Lock()
	defer f.mu.Unlock()

	var items []FavoriteItem
	for _, item := range f.Items {
		if item.Type == itemType {
			items = append(items, item)
		}
	}
	return items
}

// Contains vérifie si un élément est déjà dans les favoris
func (f *Favorites) Contains(id string, itemType FavoriteType) bool {
	f.mu.Lock()
	defer f.mu.Unlock()

	for _, item := range f.Items {
		if item.ID == id && item.Type == itemType {
			return true
		}
	}
	return false
}
