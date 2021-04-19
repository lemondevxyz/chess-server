package auth

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/toms1441/chess-server/internal/rest"
	"golang.org/x/oauth2"
)

var key = []byte("changethisimmedtialyviaconfig")
var scokie = securecookie.New(key, nil)

const (
	statesuffix = "_state"
	tokensuffix = "_token"

	timeout = time.Second * 5
)

type Config struct {
	Endpoint  oauth2.Endpoint
	Config    oauth2.Config
	MeURL     string
	ID        string
	Unmarshal UnmarshalCallback
	Logout    string
}

func NewRoutes(cfg Config) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/redirect", cfg.redirect).Methods("GET")
	r.HandleFunc("/login", cfg.login).Methods("POST")
	r.HandleFunc("/logout", cfg.logout).Methods()

	mtx.Lock()
	sliceidentify = append(sliceidentify, cfg.identify)
	mtx.Unlock()

	return r
}

func (cfg Config) redirect(w http.ResponseWriter, r *http.Request) {
	state := securecookie.GenerateRandomKey(32)

	encoded, err := scokie.Encode(cfg.ID+statesuffix, state)
	if err != nil {
		rest.RespondError(w, http.StatusInternalServerError, fmt.Errorf("securecookie.Encode: %w", err))
		return
	}

	cookie := &http.Cookie{
		Name:     cfg.ID + statesuffix,
		Value:    encoded,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		MaxAge:   60,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, cfg.Config.AuthCodeURL(string(state)), http.StatusTemporaryRedirect)
}

func (cfg Config) login(w http.ResponseWriter, r *http.Request) {
	state := r.Header.Get("state")
	if len(state) == 0 {
		rest.RespondError(w, http.StatusBadRequest, fmt.Errorf("empty state"))
		return
	}

	cookie, err := r.Cookie(cfg.ID + statesuffix)
	if err != nil {
		rest.RespondError(w, http.StatusBadRequest, fmt.Errorf("cookie error: %w", err))
		return
	}

	val := ""

	err = scokie.Decode(cfg.ID+statesuffix, cookie.Value, &val)
	if err != nil {
		rest.RespondError(w, http.StatusInternalServerError, fmt.Errorf("securecookie.Encode: %w", err))
		return
	}

	if val != state {
		rest.RespondError(w, http.StatusUnauthorized, fmt.Errorf("val != state"))
		return
	}

	code := r.Header.Get("code")
	if len(code) == 0 {
		rest.RespondError(w, http.StatusForbidden, fmt.Errorf("code is empty"))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	token, err := cfg.Config.Exchange(ctx, code)
	defer cancel()

	encoded, err := scokie.Encode(cfg.ID+tokensuffix, token)
	if err != nil {
		rest.RespondError(w, http.StatusInternalServerError, err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     cfg.ID + tokensuffix,
		Value:    encoded,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	rest.RespondJSON(w, http.StatusOK, nil)
}

func (cfg Config) logout(w http.ResponseWriter, r *http.Request) {
	token := cfg.token(r)
	if token != nil {
		rest.RespondError(w, http.StatusUnauthorized, fmt.Errorf("you're not logged in"))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	client := cfg.Config.Client(ctx, token)

	form := url.Values{}
	form.Add("token", token.AccessToken)

	resp, _ := client.Post(cfg.MeURL, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	cancel()
	resp.Body.Close()

	http.SetCookie(w, &http.Cookie{
		Name:   cfg.ID + tokensuffix,
		MaxAge: -1,
	})
}
