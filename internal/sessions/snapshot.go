package sessions

import (
	"database/sql"
	"fmt"
	"net/url"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func Snapshot(from, to string) error {
	if _, err := os.Stat(to); err == nil {
		return fmt.Errorf("%s already exists", to)
	}
	db, err := sql.Open("sqlite3", "file:"+url.PathEscape(from)+"?mode=ro&_busy_timeout=5000")
	if err != nil {
		return err
	}
	defer db.Close()
	if _, err := db.Exec("VACUUM INTO ?", to); err != nil {
		os.Remove(to)
		return err
	}
	return os.Chmod(to, 0o600)
}
