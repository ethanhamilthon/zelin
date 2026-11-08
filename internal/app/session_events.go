package app

func (s *SessionService) changed(id string, err error) error {
	if err == nil && s.notify != nil {
		s.notify("session:changed", id)
	}
	return err
}
