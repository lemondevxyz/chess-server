package google

import (
	"net/url"
)

const host = "googleapis.com"

var meurl = url.URL{
	Scheme: "https",
	Host:   "www." + host,
	Path:   "/oauth2/v1/userinfo",
}

func init() {
	q := meurl.Query()
	q.Set("alt", "json")

	meurl.RawQuery = q.Encode()
}

var logouturl = url.URL{
	Scheme: "https",
	Host:   "oauth2." + host,
	Path:   "/revoke",
}
