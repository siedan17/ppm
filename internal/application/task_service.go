package application

import "github.com/daniel/ppm/internal/domain"

type TaskService struct {
	tasks    domain.TaskRepository
	projects domain.ProjectRepository
}

func NewTaskService(tasks domain.TaskRepository, projects domain.ProjectRepository) *TaskService {
	return &TaskService{tasks: tasks, projects: projects}
}

func (s *TaskService) List(f domain.TaskFilter) ([]domain.Task, error) {
	return s.tasks.List(f)
}

func (s *TaskService) ListActive() ([]domain.Project, error) {
	return s.projects.ListActive()
}

func (s *TaskService) GetByID(id int) (*domain.Task, error) {
	return s.tasks.GetByID(id)
}

func (s *TaskService) Create(t *domain.Task) error {
	return s.tasks.Create(t)
}

func (s *TaskService) Update(t *domain.Task) error {
	return s.tasks.Update(t)
}

func (s *TaskService) UpdateStatus(id int, status string) error {
	return s.tasks.UpdateStatus(id, status)
}

func (s *TaskService) Delete(id int) error {
	return s.tasks.Delete(id)
}

func (s *TaskService) ListOverdue() ([]domain.Task, error) {
	return s.tasks.ListOverdue()
}

func (s *TaskService) ListByMeeting(meetingID int) ([]domain.Task, error) {
	return s.tasks.ListByMeeting(meetingID)
}
