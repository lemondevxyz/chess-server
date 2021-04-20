package discord

import (
	"net/url"
	"path"
	"time"

	"golang.org/x/oauth2"
)

const timeout = time.Second * 5
const apipath = "/api/v6"

var endpoint = oauth2.Endpoint{}

var meurl = url.URL{
	Scheme: "https",
	Host:   host,
	Path:   path.Join(apipath, "users", "@me"),
}

var logouturl = url.URL{
	Scheme: "https",
	Host:   host,
	Path:   path.Join(apipath, "oauth2", "token", "revoke"),
}

func init() {
	auth := url.URL{
		Scheme: "https",
		Host:   host,
		Path:   path.Join(apipath, "oauth2", "authorize"),
	}
	token := url.URL{
		Scheme: "https",
		Host:   host,
		Path:   path.Join(apipath, "oauth2", "token"),
	}

	endpoint.AuthURL = auth.String()
	endpoint.TokenURL = token.String()
}
