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

/*
func GetUser(src oauth2.TokenSource) (u User, err error) {
	if src == nil {
		return u, fmt.Errorf("src == nil")
	}

	var tok *oauth2.Token
	tok, err = src.Token()
	if err != nil {
		return u, fmt.Errorf("src.Token: %w", err)
	}

	if tok.Valid() {
		return u, fmt.Errorf("token is invalid")
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	client := oauth2.NewClient(ctx, src)

	resp, err := client.Get(meurl.String())
	if err != nil {
		return u, fmt.Errorf("client.Get: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return u, fmt.Errorf("ioutil.ReadAll: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return u, fmt.Errorf("status is not 200. status: %d, body: %s", resp.StatusCode, string(body))
	}

	err = json.Unmarshal(body, &u)
	if err != nil {
		return u, fmt.Errorf("json.Decode: %w", err)
	}

	return u, nil
}
*/
