package application

import "github.com/daniel/ppm/internal/domain"

type DashboardService struct {
	dashboard domain.DashboardRepository
	tasks     domain.TaskRepository
}

func NewDashboardService(dashboard domain.DashboardRepository, tasks domain.TaskRepository) *DashboardService {
	return &DashboardService{dashboard: dashboard, tasks: tasks}
}

func (s *DashboardService) GetDashboard() (*domain.DashboardData, error) {
	return s.dashboard.GetDashboard()
}

func (s *DashboardService) ListActiveTasks() ([]domain.Task, error) {
	return s.tasks.List(domain.TaskFilter{})
}
