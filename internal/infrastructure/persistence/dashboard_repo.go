package persistence

import (
	"context"
	"fmt"

	"github.com/daniel/ppm/internal/domain"
	"github.com/daniel/ppm/internal/infrastructure/persistence/sqlcdb"
)

var _ domain.DashboardRepository = (*DashboardRepo)(nil)

type DashboardRepo struct {
	q *sqlcdb.Queries
}

func NewDashboardRepo(q *sqlcdb.Queries) *DashboardRepo {
	return &DashboardRepo{q: q}
}

func (r *DashboardRepo) GetDashboard() (*domain.DashboardData, error) {
	data := &domain.DashboardData{}
	ctx := context.Background()

	// Projects by priority with task counts
	pRows, err := r.q.GetDashboardProjects(ctx)
	if err != nil {
		return nil, err
	}
	for _, row := range pRows {
		data.Projects = append(data.Projects, domain.Project{
			ID: int(row.ID), Name: row.Name, Priority: int(row.Priority), Status: row.Status,
			TaskCount: toInt(row.TaskCount), OverdueCount: toInt(row.OverdueCount),
		})
	}

	// Overdue tasks
	tRows, err := r.q.GetOverdueTasks(ctx)
	if err != nil {
		return nil, err
	}
	for _, row := range tRows {
		data.OverdueTasks = append(data.OverdueTasks, domain.Task{
			ID: int(row.ID), Title: row.Title, Deadline: row.Deadline,
			Status: row.Status, ProjectID: int(row.ProjectID),
			ProjectName: row.ProjectName, IsOverdue: true,
		})
	}

	// Upcoming meetings
	mRows, err := r.q.GetUpcomingMeetings(ctx)
	if err != nil {
		return nil, err
	}
	for _, row := range mRows {
		data.UpcomingMeetings = append(data.UpcomingMeetings, domain.Meeting{
			ID: int(row.ID), Title: row.Title, Date: row.Date,
			ProjectID: int(row.ProjectID), ProjectName: row.ProjectName,
		})
	}

	return data, nil
}

// toInt converts interface{} from COALESCE subqueries to int
func toInt(v interface{}) int {
	switch val := v.(type) {
	case int64:
		return int(val)
	case float64:
		return int(val)
	case int:
		return val
	default:
		n := 0
		fmt.Sscanf(fmt.Sprintf("%v", v), "%d", &n)
		return n
	}
}
