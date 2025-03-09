package spotify

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
)

// Points de terminaison pour l'API Spotify
const (
	SearchEndpoint          = "/search"
	ArtistEndpoint          = "/artists/%s"
	ArtistAlbumsEndpoint    = "/artists/%s/albums"
	AlbumEndpoint           = "/albums/%s"
	AlbumTracksEndpoint     = "/albums/%s/tracks"
	TrackEndpoint           = "/tracks/%s"
	GenresEndpoint          = "/recommendations/available-genre-seeds"
	RecommendationsEndpoint = "/recommendations"
)

// SearchOptions contient les options pour les requêtes de recherche
type SearchOptions struct {
	Query  string
	Types  []string
	Limit  int
	Offset int
	Market string
}

// BuildSearchParams construit les paramètres d'URL pour les requêtes de recherche
func BuildSearchParams(options SearchOptions) url.Values {
	params := url.Values{}
	params.Add("q", options.Query)
	params.Add("type", strings.Join(options.Types, ","))
	params.Add("limit", strconv.Itoa(options.Limit))
	params.Add("offset", strconv.Itoa(options.Offset))

	if options.Market != "" {
		params.Add("market", options.Market)
	}

	return params
}

// ParseGenresResponse analyse la réponse du point de terminaison des genres
func ParseGenresResponse(body []byte) ([]string, error) {
	var response struct {
		Genres []string `json:"genres"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response.Genres, nil
}

// ParseRecommendationsResponse analyse la réponse du point de terminaison des recommandations
func ParseRecommendationsResponse(body []byte) ([]Track, error) {
	var response struct {
		Tracks []Track `json:"tracks"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response.Tracks, nil
}

// BuildRecommendationsParams construit les paramètres d'URL pour les requêtes de recommandations
func BuildRecommendationsParams(seedGenres []string, limit int) url.Values {
	params := url.Values{}
	params.Add("seed_genres", strings.Join(seedGenres, ","))
	params.Add("limit", strconv.Itoa(limit))
	return params
}

// FormatArtistSearchParams formate les paramètres pour les recherches d'artistes
func FormatArtistSearchParams(query string, limit, offset int) url.Values {
	return BuildSearchParams(SearchOptions{
		Query:  query,
		Types:  []string{"artist"},
		Limit:  limit,
		Offset: offset,
	})
}

// FormatAlbumSearchParams formate les paramètres pour les recherches d'albums
func FormatAlbumSearchParams(query string, limit, offset int) url.Values {
	return BuildSearchParams(SearchOptions{
		Query:  query,
		Types:  []string{"album"},
		Limit:  limit,
		Offset: offset,
	})
}

// FormatTrackSearchParams formate les paramètres pour les recherches de pistes
func FormatTrackSearchParams(query string, limit, offset int) url.Values {
	return BuildSearchParams(SearchOptions{
		Query:  query,
		Types:  []string{"track"},
		Limit:  limit,
		Offset: offset,
	})
}

// FormatPaginationParams formate les paramètres pour les requêtes paginées
func FormatPaginationParams(limit, offset int) url.Values {
	params := url.Values{}
	params.Add("limit", strconv.Itoa(limit))
	params.Add("offset", strconv.Itoa(offset))
	return params
}
