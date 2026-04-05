package domain

const (
	ProjectActive    = "active"
	ProjectOnHold    = "on_hold"
	ProjectCompleted = "completed"
	ProjectArchived  = "archived"
)

var ProjectStatuses = []string{ProjectActive, ProjectOnHold, ProjectCompleted, ProjectArchived}
var Priorities = []int{1, 2, 3, 4, 5}

type Project struct {
	ID          int
	Name        string
	Priority    int
	StartDate   string
	EndDate     string
	Status      string
	GeneralInfo string
	StaticInfo  string
	DynamicInfo string
	Timeline    string
	CreatedAt   string
	UpdatedAt   string

	// Joined data (not always populated)
	People       []Person
	TaskCount    int
	OverdueCount int
}

type ProjectRepository interface {
	List() ([]Project, error)
	ListActive() ([]Project, error)
	GetByID(id int) (*Project, error)
	Create(p *Project) error
	Update(p *Project) error
	Delete(id int) error
	LinkPerson(projectID, personID int, roleInProject string) error
	UnlinkPerson(projectID, personID int) error
}
