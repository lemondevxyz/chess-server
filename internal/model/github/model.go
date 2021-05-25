package github

import (
	"strconv"

	"github.com/toms1441/chess-server/internal/model"
)

type User struct {
	// ID of the user
	ID int `json:"id"`
	// Name of the user, could be fullname or username!
	Name string `json:"login"`
	// Picture is direct URL for the picture
	Picture string `json:"avatar_url"`
}

const platform = "github"

func (u User) GetProfile() model.Profile {
	return model.Profile{
		ID:       strconv.Itoa(u.ID),
		Picture:  u.Picture,
		Username: u.Name,
		Platform: platform,
	}
}
