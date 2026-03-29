package models

type Project struct {
	ID          int
	Name        string
	Priority    int
	StartDate   string
	EndDate     string
	Status      string
	StaticInfo  string
	DynamicInfo string
	CreatedAt   string
	UpdatedAt   string

	// Joined data (not always populated)
	People       []Person
	TaskCount    int
	OverdueCount int
}
