package model

// Profile is an interface used to get information about the user.
// Having it as an interface makes it easy to swap oauth2 platforms.
type Profile interface {
	// GetPicture returns a URL of the profile picture of the user
	GetPicture() string
	// GetUsername returns a the name of the user
	GetUsername() string
	// GetPublicID returns the public id for the user, public id is primarily used by the invite system.
	GetPublicID() string
}

// AuthUser is an interface used to the authenticate user.
type AuthUser interface {
	// GetToken returns the token for the user, token is primarily used to authenticate the commands.
	GetToken() string
	// GetProfile returns the profile for the user.
	GetProfile() Profile
}
