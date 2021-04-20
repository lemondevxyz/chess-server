package model

// ProfileGetter is an interface to get profile
type ProfileGetter interface {
	GetProfile() Profile
}

// Profile is a generic struct used to get information about the user.
type Profile struct {
	ID       string `json:"id"`
	Picture  string `json:"picture"`
	Username string `json:"username"`
	Platform string `json:"platform"`
}

func (p Profile) Valid() bool {
	return len(p.ID) > 0 && len(p.Picture) > 0 && len(p.Username) > 0
}
