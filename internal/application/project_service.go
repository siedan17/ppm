package application

import "github.com/daniel/ppm/internal/domain"

type ProjectDetail struct {
	Project   *domain.Project
	Tasks     []domain.Task
	Meetings  []domain.Meeting
	AllPeople []domain.Person
}

type ProjectService struct {
	projects domain.ProjectRepository
	tasks    domain.TaskRepository
	meetings domain.MeetingRepository
	people   domain.PersonRepository
}

func NewProjectService(projects domain.ProjectRepository, tasks domain.TaskRepository, meetings domain.MeetingRepository, people domain.PersonRepository) *ProjectService {
	return &ProjectService{projects: projects, tasks: tasks, meetings: meetings, people: people}
}

func (s *ProjectService) List() ([]domain.Project, error) {
	return s.projects.List()
}

func (s *ProjectService) ListActive() ([]domain.Project, error) {
	return s.projects.ListActive()
}

func (s *ProjectService) GetDetail(id int) (*ProjectDetail, error) {
	project, err := s.projects.GetByID(id)
	if err != nil {
		return nil, err
	}
	tasks, _ := s.tasks.ListByProject(id)
	meetings, _ := s.meetings.ListByProject(id)
	allPeople, _ := s.people.List()
	return &ProjectDetail{Project: project, Tasks: tasks, Meetings: meetings, AllPeople: allPeople}, nil
}

func (s *ProjectService) GetByID(id int) (*domain.Project, error) {
	return s.projects.GetByID(id)
}

func (s *ProjectService) Create(p *domain.Project) error {
	return s.projects.Create(p)
}

func (s *ProjectService) Update(p *domain.Project) error {
	return s.projects.Update(p)
}

func (s *ProjectService) Delete(id int) error {
	return s.projects.Delete(id)
}

func (s *ProjectService) LinkPerson(projectID, personID int, roleInProject string) error {
	return s.projects.LinkPerson(projectID, personID, roleInProject)
}

func (s *ProjectService) UnlinkPerson(projectID, personID int) error {
	return s.projects.UnlinkPerson(projectID, personID)
}
