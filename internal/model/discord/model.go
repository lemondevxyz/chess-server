package discord

// Package discord provides user model, and authentication config for rest/auth.
// As well as helper function to get profile pictures, usernames and such.

import (
	"fmt"
	"net/url"
	"path"
	"strconv"

	"github.com/toms1441/chess-server/internal/model"
)

type User struct {
	// ID of the user
	ID string `json:"id"`
	// Nickname of the user(without the #). E.g. Nelly
	Nickname string `json:"username"`
	// Discriminator is the number after the username with the hashtag. E.g. #4444
	Discriminator string `json:"discriminator"` // # + the 4 numbers fater a username
	// Avatar is the hash used to get the url to the profile picture.
	Avatar string `json:"avatar"`
}

const host = "discord.com"
const platform = "discord"

var cdn = url.URL{
	Scheme: "https",
	Host:   "cdn.discordapp.com",
}

func (u User) GetPicture() string {
	pic, _ := cdn.Parse(cdn.String())
	if len(u.Avatar) > 0 {
		ext := "png"
		if u.Avatar[0] == 'a' && u.Avatar[1] == '_' {
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

func (u User) GetProfile() model.Profile {
	return model.Profile{
		ID:       u.ID,
		Picture:  u.GetPicture(),
		Username: u.GetUsername(),
		Platform: platform,
	}
}
