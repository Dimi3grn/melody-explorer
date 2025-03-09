package models

import "time"

// Album représente un album musical
type Album struct {
	ID                   string            `json:"id"`
	Name                 string            `json:"name"`
	Type                 string            `json:"album_type"`
	Artists              []Artist          `json:"artists"`
	Images               []Image           `json:"images"`
	ReleaseDate          string            `json:"release_date"`
	ReleaseDatePrecision string            `json:"release_date_precision"`
	TotalTracks          int               `json:"total_tracks"`
	Genres               []string          `json:"genres"`
	Label                string            `json:"label"`
	Popularity           int               `json:"popularity"`
	ExternalURLs         map[string]string `json:"external_urls"`
	URI                  string            `json:"uri"`
}

// AlbumPage représente une page d'albums
type AlbumPage struct {
	Items    []Album `json:"items"`
	Total    int     `json:"total"`
	Limit    int     `json:"limit"`
	Offset   int     `json:"offset"`
	Previous string  `json:"previous"`
	Next     string  `json:"next"`
}

// ParsedReleaseDate renvoie la date de sortie sous forme de time.Time
func (a *Album) ParsedReleaseDate() (time.Time, error) {
	// Le format dépend de la précision
	format := "2006-01-02"
	if a.ReleaseDatePrecision == "month" {
		format = "2006-01"
	} else if a.ReleaseDatePrecision == "year" {
		format = "2006"
	}

	return time.Parse(format, a.ReleaseDate)
}

// GetReleaseYear renvoie uniquement la composante année de la date de sortie
func (a *Album) GetReleaseYear() string {
	// Essayer d'analyser la date de sortie
	t, err := a.ParsedReleaseDate()
	if err != nil {
		// Si nous ne pouvons pas l'analyser, renvoyer les 4 premiers caractères (probablement l'année)
		if len(a.ReleaseDate) >= 4 {
			return a.ReleaseDate[:4]
		}
		return a.ReleaseDate
	}

	return t.Format("2006")
}

// PrimaryImage renvoie l'image principale de l'album ou une chaîne vide
func (a *Album) PrimaryImage() string {
	if len(a.Images) == 0 {
		return ""
	}
	return a.Images[0].URL
}

// PrimaryArtist renvoie l'artiste principal de l'album
func (a *Album) PrimaryArtist() *Artist {
	if len(a.Artists) == 0 {
		return nil
	}
	return &a.Artists[0]
}

// GetSpotifyURL renvoie l'URL Spotify pour l'album
func (a *Album) GetSpotifyURL() string {
	if url, ok := a.ExternalURLs["spotify"]; ok {
		return url
	}
	return ""
}
