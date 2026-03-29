package repository

import (
	"github.com/daniel/ppm/internal/database"
	"github.com/daniel/ppm/internal/models"
)

type DashboardRepo struct {
	db *database.DB
}

func NewDashboardRepo(db *database.DB) *DashboardRepo {
	return &DashboardRepo{db: db}
}

type DashboardData struct {
	Projects         []models.Project
	OverdueTasks     []models.Task
	UpcomingMeetings []models.Meeting
}

func (r *DashboardRepo) GetDashboard() (*DashboardData, error) {
	data := &DashboardData{}

	// Projects by priority with task counts
	rows, err := r.db.Query(`
		SELECT p.id, p.name, p.priority, p.status,
			COALESCE((SELECT COUNT(*) FROM tasks t WHERE t.project_id = p.id AND t.status NOT IN ('done','cancelled')), 0) as task_count,
			COALESCE((SELECT COUNT(*) FROM tasks t WHERE t.project_id = p.id AND t.deadline < date('now') AND t.status NOT IN ('done','cancelled')), 0) as overdue_count
		FROM projects p
		WHERE p.status = 'active'
		ORDER BY p.priority ASC, p.name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p models.Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Priority, &p.Status, &p.TaskCount, &p.OverdueCount); err != nil {
			return nil, err
		}
		data.Projects = append(data.Projects, p)
	}

	// Overdue tasks
	taskRows, err := r.db.Query(`
		SELECT t.id, t.title, t.deadline, t.status, t.project_id, p.name as project_name
		FROM tasks t
		JOIN projects p ON p.id = t.project_id
		WHERE t.deadline < date('now') AND t.status NOT IN ('done','cancelled')
		ORDER BY t.deadline ASC`)
	if err != nil {
		return nil, err
	}
	defer taskRows.Close()

	for taskRows.Next() {
		var t models.Task
		if err := taskRows.Scan(&t.ID, &t.Title, &t.Deadline, &t.Status, &t.ProjectID, &t.ProjectName); err != nil {
			return nil, err
		}
		t.IsOverdue = true
		data.OverdueTasks = append(data.OverdueTasks, t)
	}

	// Upcoming meetings (next 7 days)
	meetRows, err := r.db.Query(`
		SELECT m.id, m.title, m.date, m.project_id, p.name as project_name
		FROM meetings m
		JOIN projects p ON p.id = m.project_id
		WHERE m.date BETWEEN date('now') AND date('now', '+7 days')
		ORDER BY m.date ASC`)
	if err != nil {
		return nil, err
	}
	defer meetRows.Close()

	for meetRows.Next() {
		var m models.Meeting
		if err := meetRows.Scan(&m.ID, &m.Title, &m.Date, &m.ProjectID, &m.ProjectName); err != nil {
			return nil, err
		}
		data.UpcomingMeetings = append(data.UpcomingMeetings, m)
	}

	return data, nil
}
