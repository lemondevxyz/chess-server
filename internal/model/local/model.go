package local

// Package local provides a user model to use while debugging/testing.

import (
	"fmt"

	"github.com/toms1441/chess-server/internal/model"
)

var i int

type User struct {
	// ID of the user
	ID string `json:"id"`
	// Nickname of the user(without the #). E.g. Nelly
	Username string `json:"username"`
}

const platform = "127.0.0.1"
const picture = "https://lemondev.xyz/android-icon-192x192.png"

func NewUser() model.Profile {
	i++
	id := fmt.Sprintf("#%04d", i)
	return model.Profile{
		ID:       id,
		Picture:  picture,
		Username: id,
		Platform: platform,
	}
}
