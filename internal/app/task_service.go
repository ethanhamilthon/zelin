package app

import "zelin/internal/tasks"

type TaskService struct {
	tasks *tasks.Manager
}

func (t *TaskService) List(session string) []tasks.Info {
	if found := t.tasks.Find(session); found != nil {
		return found.List(true)
	}
	return []tasks.Info{}
}

func (t *TaskService) Output(session string, id int) (string, error) {
	return t.tasks.Session(session).Output(id, 200)
}

func (t *TaskService) Kill(session string, id int) error {
	return t.tasks.Session(session).Kill(id, false)
}
