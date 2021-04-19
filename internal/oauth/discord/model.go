package discord

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path"
	"strconv"
)

type User struct {
	// ID of the user
	ID string `json:"id"`
	// Nickname of the user(without the #). E.g. Nelly
	Nickname string `json:"username"`
	// Discriminator is the number after the username with the hashtag. E.g. #4444
	Discriminator string `json:"discriminator"` // # + the 4 numbers fater a username
	// Avatar is the hash used to get the url to the profile picture.
	Avatar string `json:"string"`
}

const host = "discord.com"
const platform = "discord"

var cdn = url.URL{
	Scheme: "https",
	Host:   fmt.Sprintf("cdn.%s", host),
}

func (u User) GetPicture() string {
	pic, _ := cdn.Parse(cdn.String())
	if len(u.Avatar) > 0 {
		ext := "png"
		if u.Avatar[0] == 'a' {
			ext = "gif"
		}

		file := fmt.Sprintf("%s.%s", u.Avatar, ext)
		pic.Path = path.Join("avatars", u.ID, file)
	} else {
		disc, _ := strconv.Atoi(u.Discriminator)

		file := fmt.Sprintf("%d.png", disc%5)
		pic.Path = path.Join("embed", "avatars", file)
	}

	return pic.String()
}

func (u User) GetUsername() string {
	return fmt.Sprintf("%s#%s", u.Nickname, u.Discriminator)
}

func (u User) GetPublicID() string {
	return u.ID
}

func (u User) GetPlatform() string {
	return platform
}

func (u User) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{
		"id":            u.ID,
		"nickname":      u.Nickname,
		"discriminator": u.Discriminator,
		"username":      u.GetUsername(),
		"avatar":        u.Avatar,
		"picture":       u.GetPicture(),
		"platform":      u.GetPlatform(),
	}

	return json.Marshal(m)
}
