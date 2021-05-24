package google

import (
	"github.com/toms1441/chess-server/internal/model"
)

type User struct {
	// ID of the user
	ID string `json:"id"`
	// Fullname of the user
	Name string `json:"name"`
	// Picture is direct URL for the picture
	Picture string `json:"picture"`
}

const platform = "google"

func (u User) GetProfile() model.Profile {
	return model.Profile{
		ID:       u.ID,
		Picture:  u.Picture,
		Username: u.Name,
		Platform: platform,
	}
}
