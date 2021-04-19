package local

import "fmt"

var i int

type User struct {
	// ID of the user
	ID string `json:"id"`
	// Nickname of the user(without the #). E.g. Nelly
	Username string `json:"username"`
}

const platform = "local"
const picture = "https://lemondev.xyz/android-icon-192x192.png"

func NewUser() User {
	i++
	id := fmt.Sprintf("#%04d", i)
	return User{
		ID:       id,
		Username: id,
	}
}

func (u User) GetPicture() string  { return picture }
func (u User) GetUsername() string { return u.Username }
func (u User) GetPublicID() string { return u.ID }
func (u User) GetPlatform() string { return platform }
