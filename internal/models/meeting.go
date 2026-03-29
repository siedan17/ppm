package models

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
