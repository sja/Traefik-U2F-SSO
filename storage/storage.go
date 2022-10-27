package storage

import (
	"database/sql"
	"fmt"
	"github.com/Tedyst/Traefik-U2F-SSO/models"
	"github.com/Tedyst/sqlitestore"
	"go.uber.org/zap"

	"github.com/koesie10/webauthn/webauthn"
	_ "github.com/mattn/go-sqlite3"
)

// ensure interface is implemented at compile time
var _ webauthn.AuthenticatorStore = (*Storage)(nil)

// Storage is needed for webauthn protocol
type Storage struct {
	logger       *zap.SugaredLogger
	db           *sql.DB
	sessionstore *sqlitestore.SqliteStore
}

func (s *Storage) CloseDb() {
	err := s.db.Close()
	if err != nil {
		s.logger.Fatalw("error closing db: ", err)
	}
}

func (s *Storage) GetSessionsStore() *sqlitestore.SqliteStore {
	return s.sessionstore
}

// AddAuthenticator is needed for webauthn protocol
func (s *Storage) AddAuthenticator(user webauthn.User, authenticator webauthn.Authenticator) error {
	logger := s.logger.With("User", user.WebAuthName())
	stmt, err := s.db.Prepare("INSERT INTO authenticators(User, ID, CredentialID, PublicKey, AAGUID, SignCount) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		logger.Error(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(user.WebAuthName(),
		authenticator.WebAuthID(),
		authenticator.WebAuthCredentialID(),
		authenticator.WebAuthPublicKey(),
		authenticator.WebAuthAAGUID(),
		authenticator.WebAuthSignCount())
	if err != nil {
		logger.Error(err)
		// return fmt.Errorf("authenticator already exists")
	}
	logger.Debugw("Added authenticator in database",
		"AuthID", authenticator.WebAuthID())
	return nil
}

// GetAuthenticator is needed for webauthn protocol
func (s *Storage) GetAuthenticator(id []byte) (webauthn.Authenticator, error) {
	logger := s.logger.With("AuthID", id)
	var au models.Authenticator
	var user string
	stmt, err := s.db.Prepare("SELECT User,ID,CredentialID,PublicKey,AAGUID,SignCount FROM authenticators WHERE ID = ?")
	if err != nil {
		logger.Error(err)
	}
	rows, err := stmt.Query(id)
	defer rows.Close()
	defer stmt.Close()
	for rows.Next() {
		err = rows.Scan(&user, &au.ID, &au.CredentialID, &au.PublicKey, &au.AAGUID, &au.SignCount)
		if err != nil {
			logger.Error(err)
		}
		logger.Debugw("Found authenticator in database",
			"User", user)
		return &au, nil
	}
	err = rows.Err()
	if err != nil {
		logger.Error(err)
	}
	logger.Debugw("Did not find authenticator in database")
	return nil, fmt.Errorf("authenticator not found")
}

// GetAuthenticators is needed for webauthn protocol
func (s *Storage) GetAuthenticators(user webauthn.User) ([]webauthn.Authenticator, error) {
	logger := s.logger.With("User", user.WebAuthName())
	var authrs []webauthn.Authenticator
	stmt, err := s.db.Prepare("SELECT ID, CredentialID, PublicKey, AAGUID, SignCount FROM authenticators WHERE User = ?")
	if err != nil {
		logger.Error(err)
	}
	rows, err := stmt.Query(user.WebAuthName())
	if err != nil {
		logger.Error(err)
		return authrs, nil
	}
	defer rows.Close()
	defer stmt.Close()
	for rows.Next() {
		var au models.Authenticator
		err = rows.Scan(&au.ID, &au.CredentialID, &au.PublicKey, &au.AAGUID, &au.SignCount)
		if err != nil {
			logger.Error(err)
		}
		logger.Debugw("Found authenticator in database",
			"AuthID", au.ID)
		authrs = append(authrs, &au)
	}
	err = rows.Err()
	if err != nil {
		logger.Error(err)
	}
	return authrs, nil
}
