package sessions

func modelKey(id string) string  { return "session-model:" + id }
func effortKey(id string) string { return "session-effort:" + id }

func (s *Store) Model(id string) (string, string, error) {
	model, err := s.Setting(modelKey(id))
	if err != nil || model == "" {
		return "", "", err
	}
	effort, err := s.Setting(effortKey(id))
	return model, effort, err
}

func (s *Store) SetModel(id, model, effort string) error {
	if err := s.SetSetting(modelKey(id), model); err != nil {
		return err
	}
	return s.SetSetting(effortKey(id), effort)
}

func (s *Store) PinModels(model, effort string) error {
	for key, value := range map[string]string{"session-model:": model, "session-effort:": effort} {
		_, err := s.db.Exec("INSERT OR IGNORE INTO settings(key,value) SELECT ?||id, ? FROM sessions WHERE messages NOT IN ('', '[]', 'null')", key, value)
		if err != nil {
			return err
		}
	}
	return nil
}
