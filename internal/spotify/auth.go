package spotify

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/spotify"
)

// Auth gère l'authentification Spotify
type Auth struct {
	Config       *oauth2.Config
	State        string
	AccessToken  string
	RefreshToken string
	Expiry       time.Time
}

// NewAuth crée une nouvelle instance Auth
func NewAuth() (*Auth, error) {
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	redirectURI := os.Getenv("REDIRECT_URI")

	if clientID == "" || clientSecret == "" || redirectURI == "" {
		return nil, fmt.Errorf("informations d'identification Spotify manquantes dans les variables d'environnement")
	}

	// Générer un état aléatoire
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	state := base64.StdEncoding.EncodeToString(b)

	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURI,
		Scopes: []string{
			"user-read-private",
			"user-read-email",
			"user-top-read",
			"user-library-read",
			"playlist-read-private",
			"playlist-read-collaborative",
		},
		Endpoint: spotify.Endpoint,
	}

	return &Auth{
		Config: config,
		State:  state,
	}, nil
}

// GetAuthURL renvoie l'URL d'autorisation Spotify
func (a *Auth) GetAuthURL() string {
	return a.Config.AuthCodeURL(a.State)
}

// Exchange échange un code d'autorisation contre des jetons
func (a *Auth) Exchange(code string) error {
	token, err := a.Config.Exchange(context.Background(), code)
	if err != nil {
		return err
	}

	a.AccessToken = token.AccessToken
	a.RefreshToken = token.RefreshToken
	a.Expiry = token.Expiry

	return nil
}

// RefreshAccessToken rafraîchit le jeton d'accès en utilisant le jeton de rafraîchissement
func (a *Auth) RefreshAccessToken() error {
	if a.RefreshToken == "" {
		return fmt.Errorf("aucun jeton de rafraîchissement disponible")
	}

	token := &oauth2.Token{
		RefreshToken: a.RefreshToken,
		Expiry:       a.Expiry,
	}

	tokenSource := a.Config.TokenSource(context.Background(), token)
	newToken, err := tokenSource.Token()
	if err != nil {
		return err
	}

	a.AccessToken = newToken.AccessToken
	// Conserver le jeton de rafraîchissement si le nouveau est vide
	if newToken.RefreshToken != "" {
		a.RefreshToken = newToken.RefreshToken
	}
	a.Expiry = newToken.Expiry

	return nil
}

// IsTokenValid vérifie si le jeton d'accès est toujours valide
func (a *Auth) IsTokenValid() bool {
	return a.AccessToken != "" && time.Now().Before(a.Expiry)
}

// EnsureValidToken s'assure que le jeton est valide, en le rafraîchissant si nécessaire
func (a *Auth) EnsureValidToken() error {
	if !a.IsTokenValid() && a.RefreshToken != "" {
		return a.RefreshAccessToken()
	}
	return nil
}

// AuthMiddleware est un middleware qui garantit l'existence d'un jeton valide
func (a *Auth) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ignorer la vérification d'authentification pour les points de terminaison liés à la connexion
		if r.URL.Path == "/login" || r.URL.Path == "/callback" || r.URL.Path == "/" {
			next.ServeHTTP(w, r)
			return
		}

		// Vérifier si nous avons un jeton valide
		if err := a.EnsureValidToken(); err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}
