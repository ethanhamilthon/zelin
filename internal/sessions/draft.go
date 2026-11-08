package sessions

import (
	"encoding/json"

	"zelin/internal/tagged"
)

func (s *Store) SaveRefs(id string, refs []tagged.Ref) error {
	if refs == nil {
		refs = []tagged.Ref{}
	}
	data, err := json.Marshal(refs)
	if err != nil {
		return err
	}
	result, err := s.db.Exec("UPDATE sessions SET refs=? WHERE id=?", string(data), id)
	return changed(result, err)
}

func (s *Store) SaveDraft(id, draft string) error {
	result, err := s.db.Exec("UPDATE sessions SET draft=?,updated=strftime('%Y-%m-%dT%H:%M:%fZ','now') WHERE id=?", draft, id)
	return changed(result, err)
}

func (s *Store) SaveAttachments(id string, attachments []string) error {
	if attachments == nil {
		attachments = []string{}
	}
	files, err := json.Marshal(attachments)
	if err != nil {
		return err
	}
	result, err := s.db.Exec("UPDATE sessions SET attachments=? WHERE id=?", string(files), id)
	return changed(result, err)
}
