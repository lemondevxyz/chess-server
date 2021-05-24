package google

import (
	"net/url"
	"path"
)

const host = "googleapis.com"

var meurl = url.URL{
	Scheme: "https",
	Host:   host,
	Path:   path.Join("oauth2", "v2", "userinfo?alt=\"json\""),
}

var logouturl = url.URL{
	Scheme: "https",
	Host:   "oauth2." + host,
	Path:   "/revoke",
}
