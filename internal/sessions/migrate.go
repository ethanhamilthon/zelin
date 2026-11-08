package sessions

import (
	"database/sql"
	"strings"
)

func migrate(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS workspaces (
 id TEXT PRIMARY KEY, name TEXT NOT NULL, icon TEXT NOT NULL DEFAULT '', dir TEXT NOT NULL, created TEXT NOT NULL)`)
	if err != nil {
		return err
	}
	for _, column := range []string{`workspace_id TEXT NOT NULL DEFAULT ''`, `refs TEXT NOT NULL DEFAULT '[]'`, `pending TEXT NOT NULL DEFAULT '[]'`} {
		_, err = db.Exec(`ALTER TABLE sessions ADD COLUMN ` + column)
		if err != nil && !strings.Contains(err.Error(), "duplicate column") {
			return err
		}
	}
	return nil
}
