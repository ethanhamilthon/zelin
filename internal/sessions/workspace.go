package sessions

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"zelin/internal/workspace"
)

const DefaultWorkspace = "default"

type Workspace struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
	Dir  string `json:"dir"`
}

func clean(name, icon, dir string) (Workspace, error) {
	name, icon = strings.TrimSpace(name), strings.TrimSpace(icon)
	if name == "" {
		return Workspace{}, fmt.Errorf("Name is required")
	}
	if len([]rune(name)) > 60 || len([]rune(icon)) > 16 {
		return Workspace{}, fmt.Errorf("Name or icon is too long")
	}
	path, err := workspace.Resolve(dir)
	return Workspace{Name: name, Icon: icon, Dir: path}, err
}

func (s *Store) EnsureDefault(dir string) error {
	_, err := s.db.Exec("INSERT OR IGNORE INTO workspaces(id,name,icon,dir,created) VALUES(?,?,?,?,?)",
		DefaultWorkspace, "Default", "", dir, time.Now().UTC().Format(time.RFC3339Nano))
	if err != nil {
		return err
	}
	_, err = s.db.Exec("UPDATE sessions SET workspace_id=? WHERE workspace_id=''", DefaultWorkspace)
	return err
}

func (s *Store) CreateWorkspace(name, icon, dir string) (Workspace, error) {
	item, err := clean(name, icon, dir)
	if err != nil {
		return item, err
	}
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return item, err
	}
	item.ID = hex.EncodeToString(bytes)
	_, err = s.db.Exec("INSERT INTO workspaces(id,name,icon,dir,created) VALUES(?,?,?,?,?)",
		item.ID, item.Name, item.Icon, item.Dir, time.Now().UTC().Format(time.RFC3339Nano))
	return item, err
}

func (s *Store) UpdateWorkspace(id, name, icon, dir string) (Workspace, error) {
	item, err := clean(name, icon, dir)
	if err != nil {
		return item, err
	}
	item.ID = id
	result, err := s.db.Exec("UPDATE workspaces SET name=?,icon=?,dir=? WHERE id=?", item.Name, item.Icon, item.Dir, id)
	return item, changed(result, err)
}

func (s *Store) Workspace(id string) (Workspace, error) {
	var item Workspace
	err := s.db.QueryRow("SELECT id,name,icon,dir FROM workspaces WHERE id=?", id).Scan(&item.ID, &item.Name, &item.Icon, &item.Dir)
	return item, err
}

func (s *Store) Workspaces() ([]Workspace, error) {
	rows, err := s.db.Query("SELECT id,name,icon,dir FROM workspaces ORDER BY created,id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []Workspace{}
	for rows.Next() {
		var item Workspace
		if err := rows.Scan(&item.ID, &item.Name, &item.Icon, &item.Dir); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}
