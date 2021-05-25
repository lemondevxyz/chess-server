package github

import (
	"net/url"
)

const host = "api.github.com"

var meurl = url.URL{
	Scheme: "https",
	Host:   host,
	Path:   "/user",
}
