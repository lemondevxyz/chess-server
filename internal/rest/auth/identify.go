package auth

import (
	"context"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/toms1441/chess-server/internal/model"
	"golang.org/x/oauth2"
)

type UnmarshalCallback func(io.ReadCloser) model.AuthUser
type IdentifyCallback func(http.ResponseWriter, *http.Request) model.AuthUser

var mtx sync.Mutex

func (cfg Config) token(req *http.Request) *oauth2.Token {
	cookie, err := req.Cookie(cfg.ID + tokensuffix)
	if err != nil {
		return nil
	}

	token := &oauth2.Token{}

	val := cookie.Value
	if err := scokie.Decode(cfg.ID+tokensuffix, val, token); err != nil {
		return nil
	}

	return token
}

// does not write anything, just returns a model.AuthUser
func (cfg Config) identify(req *http.Request) model.AuthUser {
	token := cfg.token(req)
	if token == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	client := cfg.Config.Client(ctx, token)

	resp, err := client.Get(cfg.MeURL)
	defer cancel()
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	return cfg.Unmarshal(resp.Body)
}
