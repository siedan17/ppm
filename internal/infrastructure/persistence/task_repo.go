package persistence

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/daniel/ppm/internal/domain"
	"github.com/daniel/ppm/internal/infrastructure/persistence/sqlcdb"
)

var _ domain.TaskRepository = (*TaskRepo)(nil)

type TaskRepo struct {
	q *sqlcdb.Queries
}

func NewTaskRepo(q *sqlcdb.Queries) *TaskRepo {
	return &TaskRepo{q: q}
}

func (r *TaskRepo) List(f domain.TaskFilter) ([]domain.Task, error) {
	params := sqlcdb.ListTasksParams{}

	if f.ProjectID > 0 {
		params.ProjectID = int64(f.ProjectID)
	}
	if f.Status != "" {
		params.Status = f.Status
	}
	if f.Category != "" {
		params.Category = f.Category
	}
	if f.Overdue {
		params.Overdue = 1
	}

	rows, err := r.q.ListTasks(context.Background(), params)
	if err != nil {
		return nil, err
	}

	tasks := make([]domain.Task, len(rows))
	for i, row := range rows {
		tasks[i] = mapListTaskRow(row)
	}
	return tasks, nil
}

func (r *TaskRepo) GetByID(id int) (*domain.Task, error) {
	row, err := r.q.GetTaskByID(context.Background(), int64(id))
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("task not found")
	}
	if err != nil {
		return nil, err
	}

	t := &domain.Task{
		ID: int(row.ID), ProjectID: int(row.ProjectID), Title: row.Title,
		StartDate: row.StartDate, Deadline: row.Deadline, EstimatedHours: row.EstimatedHours,
		Status: row.Status, Category: row.Category, IsExternal: row.IsExternal != 0,
		Description: row.Description, CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt,
		ProjectName: row.ProjectName, MeetingTitle: row.MeetingTitle,
	}
	if row.MeetingID.Valid {
		mid := int(row.MeetingID.Int64)
		t.MeetingID = &mid
	}
	return t, nil
}

func (r *TaskRepo) Create(t *domain.Task) error {
	params := sqlcdb.CreateTaskParams{
		ProjectID: int64(t.ProjectID), Title: t.Title,
		StartDate: t.StartDate, Deadline: t.Deadline, EstimatedHours: t.EstimatedHours,
		Status: t.Status, Category: t.Category, IsExternal: boolToInt64(t.IsExternal),
		Description: t.Description,
	}
	if t.MeetingID != nil {
		params.MeetingID = sql.NullInt64{Int64: int64(*t.MeetingID), Valid: true}
	}

	id, err := r.q.CreateTask(context.Background(), params)
	if err != nil {
		return err
	}
	t.ID = int(id)
	return nil
}

func (r *TaskRepo) Update(t *domain.Task) error {
	params := sqlcdb.UpdateTaskParams{
		ProjectID: int64(t.ProjectID), Title: t.Title,
		StartDate: t.StartDate, Deadline: t.Deadline, EstimatedHours: t.EstimatedHours,
		Status: t.Status, Category: t.Category, IsExternal: boolToInt64(t.IsExternal),
		Description: t.Description, ID: int64(t.ID),
	}
	if t.MeetingID != nil {
		params.MeetingID = sql.NullInt64{Int64: int64(*t.MeetingID), Valid: true}
	}
	return r.q.UpdateTask(context.Background(), params)
}

func (r *TaskRepo) UpdateStatus(id int, status string) error {
	return r.q.UpdateTaskStatus(context.Background(), sqlcdb.UpdateTaskStatusParams{
		Status: status, ID: int64(id),
	})
}

func (r *TaskRepo) Delete(id int) error {
	return r.q.DeleteTask(context.Background(), int64(id))
}

func (r *TaskRepo) ListByProject(projectID int) ([]domain.Task, error) {
	return r.List(domain.TaskFilter{ProjectID: projectID})
}

func (r *TaskRepo) ListByMeeting(meetingID int) ([]domain.Task, error) {
	rows, err := r.q.ListTasksByMeeting(context.Background(), sql.NullInt64{Int64: int64(meetingID), Valid: true})
	if err != nil {
		return nil, err
	}
	tasks := make([]domain.Task, len(rows))
	for i, row := range rows {
		tasks[i] = domain.Task{
			ID: int(row.ID), Title: row.Title, Deadline: row.Deadline,
			Status: row.Status, Category: row.Category,
		}
	}
	return tasks, nil
}

func (r *TaskRepo) ListOverdue() ([]domain.Task, error) {
	return r.List(domain.TaskFilter{Overdue: true})
}

func (r *TaskRepo) AddDependency(taskID, dependsOnID int) error {
	return r.q.AddTaskDependency(context.Background(), sqlcdb.AddTaskDependencyParams{
		TaskID: int64(taskID), DependsOnID: int64(dependsOnID),
	})
}

func (r *TaskRepo) RemoveDependency(taskID, dependsOnID int) error {
	return r.q.RemoveTaskDependency(context.Background(), sqlcdb.RemoveTaskDependencyParams{
		TaskID: int64(taskID), DependsOnID: int64(dependsOnID),
	})
}

func (r *TaskRepo) GetDependencies(taskID int) ([]domain.Task, error) {
	rows, err := r.q.GetTaskDependencies(context.Background(), int64(taskID))
	if err != nil {
		return nil, err
	}
	tasks := make([]domain.Task, len(rows))
	for i, row := range rows {
		tasks[i] = domain.Task{ID: int(row.ID), Title: row.Title, Status: row.Status}
	}
	return tasks, nil
}

func mapListTaskRow(row sqlcdb.ListTasksRow) domain.Task {
	t := domain.Task{
		ID: int(row.ID), ProjectID: int(row.ProjectID), Title: row.Title,
		StartDate: row.StartDate, Deadline: row.Deadline, EstimatedHours: row.EstimatedHours,
		Status: row.Status, Category: row.Category, IsExternal: row.IsExternal != 0,
		Description: row.Description, CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt,
		ProjectName: row.ProjectName, MeetingTitle: row.MeetingTitle,
		IsOverdue: row.IsOverdue != 0,
	}
	if row.MeetingID.Valid {
		mid := int(row.MeetingID.Int64)
		t.MeetingID = &mid
	}
	return t
}

func boolToInt64(b bool) int64 {
	if b {
		return 1
	}
	return 0
}
