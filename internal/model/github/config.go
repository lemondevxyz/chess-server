package github

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/toms1441/chess-server/internal/model"
	"github.com/toms1441/chess-server/internal/rest/auth"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

func NewAuthConfig(cfg model.OAuth2Config) auth.Config {
	return auth.Config{
		Config: oauth2.Config{
			Endpoint:     github.Endpoint,
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURL:  cfg.Redirect,
			Scopes:       []string{},
		},
		MeURL:  meurl.String(),
		ID:     "github",
		Logout: "",
		Unmarshal: func(reader io.ReadCloser) *model.Profile {
			defer reader.Close()

			/*
				body, err := ioutil.ReadAll(reader)
				fmt.Println(string(body), err)
			*/

			decode := json.NewDecoder(reader)
			user := &User{}
			if err := decode.Decode(user); err != nil {
				fmt.Println(err, user)
				return nil
			}

			pro := user.GetProfile()
			return &pro
		},
	}
}
