package config

import (
	"fmt"
	"os"
	"path/filepath"
)

type Db struct {
	SqliteFile string `mapstructure:"sqlite_file"`
}

var defaultDB = &Db{SqliteFile: "storage/database"}

func (db Db) Validate() error {
	// Check if the sqlite db path exists, db files are created automatically
	dir, _ := filepath.Abs(filepath.Dir(db.SqliteFile))
	if _, err := os.Stat(dir); err != nil {
		return fmt.Errorf("config error: path to db does not exist: %v", db.SqliteFile)
	}
	return nil
}
