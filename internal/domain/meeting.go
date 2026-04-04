package domain

const (
	MeetingInternal = "internal"
	MeetingExternal = "external"
)

var MeetingTypes = []string{MeetingInternal, MeetingExternal}

type Meeting struct {
	ID          int
	ProjectID   int
	Date        string
	MeetingType string
	Title       string
	Notes       string
	CreatedAt   string
	UpdatedAt   string

	// Joined data
	ProjectName  string
	Participants []Person
	Tasks        []Task
	NotesHTML    string
}

type MeetingFilter struct {
	ProjectID int
	DateFrom  string
	DateTo    string
}

type MeetingRepository interface {
	List(f MeetingFilter) ([]Meeting, error)
	GetByID(id int) (*Meeting, error)
	Create(m *Meeting) error
	Update(m *Meeting) error
	Delete(id int) error
	AddParticipant(meetingID, personID int) error
	RemoveParticipant(meetingID, personID int) error
	ListByProject(projectID int) ([]Meeting, error)
}
