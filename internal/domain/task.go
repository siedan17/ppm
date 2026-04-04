package domain

const (
	TaskTodo       = "todo"
	TaskInProgress = "in_progress"
	TaskBlocked    = "blocked"
	TaskDone       = "done"
	TaskCancelled  = "cancelled"
)

const (
	CategoryProgramming     = "programming"
	CategoryDataEngineering = "data_engineering"
	CategorySpecification   = "specification"
	CategoryDesign          = "design"
	CategoryCommunication   = "communication"
	CategoryOther           = "other"
)

var TaskStatuses = []string{TaskTodo, TaskInProgress, TaskBlocked, TaskDone, TaskCancelled}
var TaskCategories = []string{CategoryProgramming, CategoryDataEngineering, CategorySpecification, CategoryDesign, CategoryCommunication, CategoryOther}

type Task struct {
	ID             int
	ProjectID      int
	MeetingID      *int
	Title          string
	StartDate      string
	Deadline       string
	EstimatedHours float64
	Status         string
	Category       string
	IsExternal     bool
	Description    string
	CreatedAt      string
	UpdatedAt      string

	// Joined data
	ProjectName  string
	MeetingTitle string
	IsOverdue    bool
	DependsOn    []Task
	Blockers     []Task
}

type TaskFilter struct {
	ProjectID int
	Status    string
	Category  string
	Overdue   bool
}

type TaskRepository interface {
	List(f TaskFilter) ([]Task, error)
	GetByID(id int) (*Task, error)
	Create(t *Task) error
	Update(t *Task) error
	UpdateStatus(id int, status string) error
	Delete(id int) error
	ListByProject(projectID int) ([]Task, error)
	ListByMeeting(meetingID int) ([]Task, error)
	ListOverdue() ([]Task, error)
	AddDependency(taskID, dependsOnID int) error
	RemoveDependency(taskID, dependsOnID int) error
	GetDependencies(taskID int) ([]Task, error)
}
