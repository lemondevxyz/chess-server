package auth

import (
	"context"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/toms1441/chess-server/internal/model"
	"github.com/toms1441/chess-server/internal/model/local"
	"golang.org/x/oauth2"
)

type UnmarshalCallback func(io.ReadCloser) *model.Profile
type IdentifyCallback func(*http.Request) *model.Profile

var mtx sync.Mutex
var sliceidentify = []IdentifyCallback{}

func Identify(r *http.Request) *model.Profile {
	for _, identify := range sliceidentify {
		authuser := identify(r)
		if authuser != nil {
			return authuser
		}
	}

	if len(sliceidentify) == 0 {
		user := local.NewUser()
		return &user
	}

	return nil
}

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
func (cfg Config) identify(req *http.Request) *model.Profile {
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
