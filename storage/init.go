package storage

import (
	"database/sql"
	"fmt"
	. "github.com/Tedyst/Traefik-U2F-SSO/config"
	"github.com/Tedyst/sqlitestore"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func InitStorage(logger *zap.SugaredLogger) (*Storage, error) {
	dsn := fmt.Sprintf("%v?journal_mode=WAL", viper.GetString(ConfSqliteFile))
	logger.Debugf("open sqlite 3 at '%v'", dsn)
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	// Test storage
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	key := []byte(viper.GetString(ConfSessionKey))
	sessionsStore, err := createSessionsStore(db, key)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		return nil, err
	}

	if err := createSchema(db); err != nil {
		return nil, err
	}

	storage := Storage{
		logger:       logger,
		db:           db,
		sessionstore: sessionsStore,
	}

	return &storage, nil
}

func createSessionsStore(db *sql.DB, key []byte) (*sqlitestore.SqliteStore, error) {
	sessionsstore, err := sqlitestore.NewSqliteStoreFromConnection(db, "sessions", "/", 360000, key)
	if err != nil {
		return nil, fmt.Errorf("could not init DB: %w", err)
	}
	sessionsstore.Options.Domain = viper.GetString(ConfDomain)
	sessionsstore.Options.Secure = true
	sessionsstore.Options.HttpOnly = false
	return sessionsstore, nil
}

func createSchema(db *sql.DB) error {
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS users (
	    id integer not null primary key, name text
	);
	CREATE TABLE IF NOT EXISTS authenticators (
		User TEXT,
		ID BLOB UNIQUE,
		CredentialID BLOB,
		PublicKey BLOB,
		AAGUID BLOB,
		SignCount INTEGER
	);`
	_, err := db.Exec(sqlStmt)
	return err
}
