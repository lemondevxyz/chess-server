package google

import (
	"encoding/json"
	"io"

	"github.com/toms1441/chess-server/internal/model"
	"github.com/toms1441/chess-server/internal/rest/auth"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Config struct {
	ClientID     string `validate:"required,len=18" mapstructure:"client_id"`     // Client ID for google
	ClientSecret string `validate:"required,len=32" mapstructure:"client_secret"` // Client Secret for google
	Redirect     string `validate:"required" mapstructure:"redirect"`
}

func NewAuthConfig(cfg Config) auth.Config {
	return auth.Config{
		Config: oauth2.Config{
			Endpoint:     google.Endpoint,
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURL:  cfg.Redirect,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.profile",
			},
		},
		MeURL:  meurl.String(),
		ID:     "google",
		Logout: logouturl.String(),
		Unmarshal: func(reader io.ReadCloser) *model.Profile {
			defer reader.Close()

			decode := json.NewDecoder(reader)
			user := &User{}

			if err := decode.Decode(user); err != nil {
				return nil
			}

			pro := user.GetProfile()
			return &pro
		},
	}
}
