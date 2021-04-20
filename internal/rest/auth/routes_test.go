package auth

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"golang.org/x/oauth2"
)

var gState = ""
var cfg = Config{
	ID: "example",
	Config: oauth2.Config{
		ClientID:     "",
		ClientSecret: "ayo",
		RedirectURL:  "http://localhost:9999",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "example.com/auth",
			TokenURL: "example.com/token",
		},
		Scopes: []string{"identify"},
	},
}

func init() {
}

func TestRouteRedirect(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/redirect", nil)
	w := httptest.NewRecorder()

	cfg.redirect(w, req)
	head := w.Header()

	for _, v := range w.Result().Cookies() {
		if v.Name == "example_state" {
			gState = v.Value
			break
		}
	}

	if len(gState) == 0 {
		t.Fatalf("redirect does not set state cookie")
	}

	if len(head.Get("Location")) == 0 {
		t.Fatalf("redirect does not redirect..")
	}
	if len(head.Get("Set-Cookie")) == 0 {
		t.Fatalf("redirect does not set the cookie. invalid cookie")
	}
}

func TestRouteLogin(t *testing.T) {

	rou := cfg.login
	{ // should fail no state present
		req := httptest.NewRequest("GET", "http://example.com/login", nil)
		w := httptest.NewRecorder()

		rou(w, req)
		if w.Result().StatusCode != http.StatusBadRequest {
			t.Fatalf("Result code: %d", w.Result().StatusCode)
		}
	}

	{ // should fail state doesn't match cookie
		req := httptest.NewRequest("GET", "http://example.com/login", nil)
		req.Header.Add("state", "lamo")
		w := httptest.NewRecorder()

		rou(w, req)
		if w.Result().StatusCode != http.StatusBadRequest {
			t.Fatalf("Result code: %d", w.Result().StatusCode)
		}
	}
	{ // should fail cause no code
		earl := url.URL{
			Scheme: "http",
			Host:   "example.com",
			Path:   "/login",
		}
		q := earl.Query()
		q.Add("state", gState)

		earl.RawQuery = q.Encode()

		req := httptest.NewRequest("GET", earl.String(), nil)
		req.AddCookie(&http.Cookie{
			Name:   "example_state",
			Value:  gState,
			MaxAge: 60,
		})

		req.Header.Add("state", gState)
		w := httptest.NewRecorder()

		rou(w, req)

		if w.Result().StatusCode != http.StatusUnauthorized {
			t.Logf("body: %s", w.Body.String())
			t.Fatalf("code: %d", w.Result().StatusCode)
		}
	}
}
