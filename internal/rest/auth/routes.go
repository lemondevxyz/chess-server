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
	"github.com/thanhpk/randstr"
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
	Config    oauth2.Config
	MeURL     string
	ID        string
	Unmarshal UnmarshalCallback
	Logout    string
}

func AddRoutes(cfg Config, r *mux.Router) {
	r.HandleFunc("/redirect", cfg.redirect).Methods("GET")
	r.HandleFunc("/private", cfg.private).Methods("GET")
	r.HandleFunc("/login", cfg.login).Methods("GET")
	r.HandleFunc("/logout", cfg.logout).Methods("POST")

	mtx.Lock()
	sliceidentify = append(sliceidentify, cfg.identify)
	mtx.Unlock()
}

func (cfg Config) private(w http.ResponseWriter, r *http.Request) {
	if cfg.identify(r) == nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("not logged in"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(""))
}

func (cfg Config) redirect(w http.ResponseWriter, r *http.Request) {
	if cfg.identify(r) != nil {
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte(""))
		return
	}

	state := randstr.String(32)

	encoded, err := scokie.Encode(cfg.ID+statesuffix, state)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("securecookie.Encode: %s", err.Error())))
		return
	}

	cookie := &http.Cookie{
		Name:     cfg.ID + statesuffix,
		Value:    encoded,
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
		MaxAge:   60,
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, cfg.Config.AuthCodeURL(string(state)), http.StatusTemporaryRedirect)
}

func (cfg Config) login(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	if len(state) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("state is empty"))
		return
	}

	cookie, err := r.Cookie(cfg.ID + statesuffix)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("cookie error: " + err.Error()))
		return
	}

	val := ""

	err = scokie.Decode(cfg.ID+statesuffix, cookie.Value, &val)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("securecookie.Encode: " + err.Error()))
		return
	}

	if val != state {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("val != state"))
		return
	}

	code := r.URL.Query().Get("code")
	if len(code) == 0 {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("code is empty"))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	token, err := cfg.Config.Exchange(ctx, code)
	defer cancel()

	encoded, err := scokie.Encode(cfg.ID+tokensuffix, token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("encode: " + err.Error()))

		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  cfg.ID + tokensuffix,
		Value: encoded,
		Path:  "/",
		// Secure:   true,
		HttpOnly: true,
		// SameSite: http.SameSiteStrictMode,
		MaxAge: (3600 * 24) * 14,
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(""))
}

func (cfg Config) logout(w http.ResponseWriter, r *http.Request) {
	token := cfg.token(r)
	if token != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("not logged in"))
		return
	}

	if len(cfg.MeURL) > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		client := cfg.Config.Client(ctx, token)

		form := url.Values{}
		form.Add("token", token.AccessToken)

		resp, _ := client.Post(cfg.MeURL, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
		cancel()
		resp.Body.Close()
	}

	http.SetCookie(w, &http.Cookie{
		Name:   cfg.ID + tokensuffix,
		MaxAge: -1,
	})
}
