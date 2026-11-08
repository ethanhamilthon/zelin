package sessions

import "fmt"

func (s *Store) DeleteWorkspace(id string) error {
	if id == DefaultWorkspace {
		return fmt.Errorf("The default workspace cannot be deleted")
	}
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	result, err := tx.Exec("DELETE FROM workspaces WHERE id=?", id)
	if err := changed(result, err); err != nil {
		return err
	}
	if _, err := tx.Exec("UPDATE sessions SET workspace_id=? WHERE workspace_id=?", DefaultWorkspace, id); err != nil {
		return err
	}
	return tx.Commit()
}
