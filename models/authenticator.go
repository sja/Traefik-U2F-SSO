package models

// Authenticator is needed for webauthn protocol
type Authenticator struct {
	ID           []byte
	CredentialID []byte
	PublicKey    []byte
	AAGUID       []byte
	SignCount    uint32
}

// WebAuthID is needed for webauthn protocol
func (a *Authenticator) WebAuthID() []byte {
	return a.ID
}

// WebAuthCredentialID is needed for webauthn protocol
func (a *Authenticator) WebAuthCredentialID() []byte {
	return a.CredentialID
}

// WebAuthPublicKey is needed for webauthn protocol
func (a *Authenticator) WebAuthPublicKey() []byte {
	return a.PublicKey
}

// WebAuthAAGUID is needed for webauthn protocol
func (a *Authenticator) WebAuthAAGUID() []byte {
	return a.AAGUID
}

// WebAuthSignCount is needed for webauthn protocol
func (a *Authenticator) WebAuthSignCount() uint32 {
	return a.SignCount
}
