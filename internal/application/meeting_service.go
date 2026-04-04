package application

import "github.com/daniel/ppm/internal/domain"

// MarkdownRenderer abstracts markdown rendering so the application layer
// doesn't depend on the render package directly.
type MarkdownRenderer interface {
	RenderMarkdown(s string) string
}

type MeetingDetail struct {
	Meeting   *domain.Meeting
	Tasks     []domain.Task
	AllPeople []domain.Person
}

type MeetingService struct {
	meetings domain.MeetingRepository
	projects domain.ProjectRepository
	people   domain.PersonRepository
	tasks    domain.TaskRepository
	md       MarkdownRenderer
}

func NewMeetingService(meetings domain.MeetingRepository, projects domain.ProjectRepository, people domain.PersonRepository, tasks domain.TaskRepository, md MarkdownRenderer) *MeetingService {
	return &MeetingService{meetings: meetings, projects: projects, people: people, tasks: tasks, md: md}
}

func (s *MeetingService) List(f domain.MeetingFilter) ([]domain.Meeting, error) {
	return s.meetings.List(f)
}

func (s *MeetingService) ListActive() ([]domain.Project, error) {
	return s.projects.ListActive()
}

func (s *MeetingService) GetDetail(id int) (*MeetingDetail, error) {
	meeting, err := s.meetings.GetByID(id)
	if err != nil {
		return nil, err
	}
	meeting.NotesHTML = s.md.RenderMarkdown(meeting.Notes)
	tasks, _ := s.tasks.ListByMeeting(id)
	allPeople, _ := s.people.List()
	return &MeetingDetail{Meeting: meeting, Tasks: tasks, AllPeople: allPeople}, nil
}

func (s *MeetingService) GetByID(id int) (*domain.Meeting, error) {
	return s.meetings.GetByID(id)
}

func (s *MeetingService) Create(m *domain.Meeting) error {
	return s.meetings.Create(m)
}

func (s *MeetingService) Update(m *domain.Meeting) error {
	return s.meetings.Update(m)
}

func (s *MeetingService) Delete(id int) error {
	return s.meetings.Delete(id)
}

func (s *MeetingService) AddParticipant(meetingID, personID int) error {
	return s.meetings.AddParticipant(meetingID, personID)
}

func (s *MeetingService) RemoveParticipant(meetingID, personID int) error {
	return s.meetings.RemoveParticipant(meetingID, personID)
}
