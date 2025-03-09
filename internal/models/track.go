package models

import (
	"fmt"
	"time"
)

// Track représente une piste musicale
type Track struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Type         string            `json:"type"`
	Album        Album             `json:"album"`
	Artists      []Artist          `json:"artists"`
	DiscNumber   int               `json:"disc_number"`
	TrackNumber  int               `json:"track_number"`
	DurationMs   int               `json:"duration_ms"`
	Explicit     bool              `json:"explicit"`
	IsPlayable   bool              `json:"is_playable"`
	Popularity   int               `json:"popularity"`
	PreviewURL   string            `json:"preview_url"`
	ExternalURLs map[string]string `json:"external_urls"`
	URI          string            `json:"uri"`
}

// TrackPage représente une page de pistes
type TrackPage struct {
	Items    []Track `json:"items"`
	Total    int     `json:"total"`
	Limit    int     `json:"limit"`
	Offset   int     `json:"offset"`
	Previous string  `json:"previous"`
	Next     string  `json:"next"`
}

// PrimaryImage renvoie l'image principale de la piste à partir de son album
func (t *Track) PrimaryImage() string {
	if len(t.Album.Images) == 0 {
		return ""
	}
	return t.Album.Images[0].URL
}

// GetSpotifyURL renvoie l'URL Spotify pour la piste
func (t *Track) GetSpotifyURL() string {
	if url, ok := t.ExternalURLs["spotify"]; ok {
		return url
	}
	return ""
}

// PrimaryArtist renvoie l'artiste principal de la piste
func (t *Track) PrimaryArtist() *Artist {
	if len(t.Artists) == 0 {
		return nil
	}
	return &t.Artists[0]
}

// ArtistNames renvoie une liste des artistes de la piste séparés par des virgules
func (t *Track) ArtistNames() string {
	if len(t.Artists) == 0 {
		return "Unknown Artist"
	}

	result := t.Artists[0].Name
	for i := 1; i < len(t.Artists); i++ {
		result += ", " + t.Artists[i].Name
	}

	return result
}

// FormattedDuration renvoie la durée de la piste au format mm:ss
func (t *Track) FormattedDuration() string {
	duration := time.Duration(t.DurationMs) * time.Millisecond
	minutes := int(duration.Minutes())
	seconds := int(duration.Seconds()) % 60
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}

// HasPreview indique si la piste a une URL d'aperçu
func (t *Track) HasPreview() bool {
	return t.PreviewURL != ""
}

// DisplayName renvoie le nom de la piste ou "Unknown Track" si vide
func (t *Track) DisplayName() string {
	if t.Name == "" {
		return "Unknown Track"
	}
	return t.Name
}

// PopularityClass renvoie une classe CSS basée sur la popularité de la piste
func (t *Track) PopularityClass() string {
	if t.Popularity >= 80 {
		return "popularity-high"
	} else if t.Popularity >= 50 {
		return "popularity-medium"
	}
	return "popularity-low"
}

// ExplicitTag renvoie une étiquette si la piste contient du contenu explicite
func (t *Track) ExplicitTag() string {
	if t.Explicit {
		return "Explicit"
	}
	return ""
}
