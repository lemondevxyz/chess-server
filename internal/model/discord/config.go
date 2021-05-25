package discord

import (
	"encoding/json"
	"io"

	"github.com/toms1441/chess-server/internal/model"
	"github.com/toms1441/chess-server/internal/rest/auth"
	"golang.org/x/oauth2"
)

func NewAuthConfig(cfg model.OAuth2Config) auth.Config {
	return auth.Config{
		Config: oauth2.Config{
			Endpoint:     endpoint, // global
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURL:  cfg.Redirect,
			Scopes: []string{
				"identify",
			},
		},
		MeURL:  meurl.String(),
		ID:     "discord",
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
