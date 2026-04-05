package domain

type DashboardData struct {
	Projects         []Project
	ActiveTasks      []Task
	UpcomingMeetings []Meeting
}

type DashboardRepository interface {
	GetDashboard() (*DashboardData, error)
}
