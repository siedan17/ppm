package models

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
	ProjectName string
	MeetingTitle string
	IsOverdue   bool
	DependsOn   []Task
	Blockers    []Task
}
