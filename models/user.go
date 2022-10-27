package models

// User is needed for json reply when logged in.
type User struct {
	Name string `json:"name"`
}

// WebAuthID is needed for webauthn protocol
func (u *User) WebAuthID() []byte {
	return []byte(u.Name)
}

// WebAuthName is needed for webauthn protocol
func (u *User) WebAuthName() string {
	return u.Name
}

// WebAuthDisplayName is needed for webauthn protocol
func (u *User) WebAuthDisplayName() string {
	return u.Name
}
