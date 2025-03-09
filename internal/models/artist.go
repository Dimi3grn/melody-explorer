package models

// Artist représente un artiste musical
type Artist struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Type         string            `json:"type"`
	Images       []Image           `json:"images"`
	Genres       []string          `json:"genres"`
	Popularity   int               `json:"popularity"`
	Followers    Followers         `json:"followers"`
	ExternalURLs map[string]string `json:"external_urls"`
	URI          string            `json:"uri"`
}

// Image représente une image
type Image struct {
	URL    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

// Followers représente les informations sur les abonnés
type Followers struct {
	Total int `json:"total"`
}

// ArtistPage représente une page d'artistes
type ArtistPage struct {
	Items    []Artist `json:"items"`
	Total    int      `json:"total"`
	Limit    int      `json:"limit"`
	Offset   int      `json:"offset"`
	Previous string   `json:"previous"`
	Next     string   `json:"next"`
}

// PrimaryImage renvoie l'image principale de l'artiste ou une chaîne vide
func (a *Artist) PrimaryImage() string {
	if len(a.Images) == 0 {
		return ""
	}
	return a.Images[0].URL
}

// GetSpotifyURL renvoie l'URL Spotify pour l'artiste
func (a *Artist) GetSpotifyURL() string {
	if url, ok := a.ExternalURLs["spotify"]; ok {
		return url
	}
	return ""
}

// GetGenreString renvoie une liste des genres de l'artiste séparés par des virgules
func (a *Artist) GetGenreString() string {
	if len(a.Genres) == 0 {
		return "Unknown"
	}

	result := a.Genres[0]
	for i := 1; i < len(a.Genres); i++ {
		result += ", " + a.Genres[i]
	}

	return result
}

// DisplayName renvoie le nom de l'artiste ou "Unknown Artist" si vide
func (a *Artist) DisplayName() string {
	if a.Name == "" {
		return "Unknown Artist"
	}
	return a.Name
}

// PopularityClass renvoie une classe CSS basée sur la popularité de l'artiste
func (a *Artist) PopularityClass() string {
	if a.Popularity >= 80 {
		return "popularity-high"
	} else if a.Popularity >= 50 {
		return "popularity-medium"
	}
	return "popularity-low"
}

// GetFollowersFormatted renvoie une chaîne formatée du nombre d'abonnés
func (a *Artist) GetFollowersFormatted() string {
	if a.Followers.Total > 1000000 {
		return formatFloat(float64(a.Followers.Total)/1000000) + "M"
	} else if a.Followers.Total > 1000 {
		return formatFloat(float64(a.Followers.Total)/1000) + "K"
	}
	return formatInt(a.Followers.Total)
}

// formatFloat formate un nombre à virgule flottante en chaîne avec jusqu'à 1 décimale
func formatFloat(num float64) string {
	if num == float64(int(num)) {
		return formatInt(int(num))
	}
	return formatInt(int(num*10)/10) + "." + formatInt(int(num*10)%10)
}

// formatInt formate un entier avec des virgules
func formatInt(num int) string {
	str := ""
	for i := 0; num > 0; i++ {
		if i > 0 && i%3 == 0 {
			str = "," + str
		}
		str = string('0'+num%10) + str
		num /= 10
	}
	if str == "" {
		return "0"
	}
	return str
}
