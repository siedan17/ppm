package domain

type DashboardData struct {
	Projects         []Project
	OverdueTasks     []Task
	UpcomingMeetings []Meeting
}

type DashboardRepository interface {
	GetDashboard() (*DashboardData, error)
}
