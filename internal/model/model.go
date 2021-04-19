package model

import "encoding/json"

// Profile is an interface used to get information about the user.
// Having it as an interface makes it easy to swap oauth2 platforms.
type Profile interface {
	// GetPicture returns a URL of the profile picture of the user
	GetPicture() string
	// GetUsername returns a the name of the user
	GetUsername() string
	// GetPublicID returns the public id for the user, public id is primarily used by the invite system.
	GetPublicID() string
	// GetPlatform returns the platform of that user, in lower case. E.g. "google", "facebook"
	GetPlatform() string

	// JSON output should contain the following
	// - username
	// - id
	// - platform: the name of the platform
	// - picture: the url to the profile picture
	// any other fields could be inserted aswell
	json.Marshaler
}
