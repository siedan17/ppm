package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/daniel/ppm/internal/database"
	"github.com/daniel/ppm/internal/models"
)

type TaskRepo struct {
	db *database.DB
}

func NewTaskRepo(db *database.DB) *TaskRepo {
	return &TaskRepo{db: db}
}

type TaskFilter struct {
	ProjectID int
	Status    string
	Category  string
	Overdue   bool
}

func (r *TaskRepo) List(f TaskFilter) ([]models.Task, error) {
	query := `
		SELECT t.id, t.project_id, t.meeting_id, t.title, t.start_date, t.deadline,
			t.estimated_hours, t.status, t.category, t.is_external, t.description,
			t.created_at, t.updated_at, p.name as project_name,
			COALESCE(m.title, '') as meeting_title,
			CASE WHEN t.deadline < date('now') AND t.status NOT IN ('done','cancelled') THEN 1 ELSE 0 END as is_overdue
		FROM tasks t
		JOIN projects p ON p.id = t.project_id
		LEFT JOIN meetings m ON m.id = t.meeting_id
		WHERE 1=1`

	var args []any

	if f.ProjectID > 0 {
		query += " AND t.project_id = ?"
		args = append(args, f.ProjectID)
	}
	if f.Status != "" {
		query += " AND t.status = ?"
		args = append(args, f.Status)
	}
	if f.Category != "" {
		query += " AND t.category = ?"
		args = append(args, f.Category)
	}
	if f.Overdue {
		query += " AND t.deadline < date('now') AND t.status NOT IN ('done','cancelled')"
	}

	query += " ORDER BY t.deadline ASC, t.priority_sort ASC"
	// priority_sort doesn't exist, sort by deadline and status
	query = strings.Replace(query, " ORDER BY t.deadline ASC, t.priority_sort ASC", " ORDER BY t.deadline ASC, t.status ASC", 1)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		var meetingID sql.NullInt64
		var isOverdue int
		if err := rows.Scan(&t.ID, &t.ProjectID, &meetingID, &t.Title, &t.StartDate, &t.Deadline,
			&t.EstimatedHours, &t.Status, &t.Category, &t.IsExternal, &t.Description,
			&t.CreatedAt, &t.UpdatedAt, &t.ProjectName, &t.MeetingTitle, &isOverdue); err != nil {
			return nil, err
		}
		if meetingID.Valid {
			id := int(meetingID.Int64)
			t.MeetingID = &id
		}
		t.IsOverdue = isOverdue == 1
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (r *TaskRepo) GetByID(id int) (*models.Task, error) {
	var t models.Task
	var meetingID sql.NullInt64
	err := r.db.QueryRow(`
		SELECT t.id, t.project_id, t.meeting_id, t.title, t.start_date, t.deadline,
			t.estimated_hours, t.status, t.category, t.is_external, t.description,
			t.created_at, t.updated_at, p.name as project_name,
			COALESCE(m.title, '') as meeting_title
		FROM tasks t
		JOIN projects p ON p.id = t.project_id
		LEFT JOIN meetings m ON m.id = t.meeting_id
		WHERE t.id = ?`, id).
		Scan(&t.ID, &t.ProjectID, &meetingID, &t.Title, &t.StartDate, &t.Deadline,
			&t.EstimatedHours, &t.Status, &t.Category, &t.IsExternal, &t.Description,
			&t.CreatedAt, &t.UpdatedAt, &t.ProjectName, &t.MeetingTitle)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("task not found")
	}
	if err != nil {
		return nil, err
	}
	if meetingID.Valid {
		id := int(meetingID.Int64)
		t.MeetingID = &id
	}
	return &t, nil
}

func (r *TaskRepo) Create(t *models.Task) error {
	var meetingID any
	if t.MeetingID != nil {
		meetingID = *t.MeetingID
	}
	result, err := r.db.Exec(`
		INSERT INTO tasks (project_id, meeting_id, title, start_date, deadline, estimated_hours, status, category, is_external, description)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		t.ProjectID, meetingID, t.Title, t.StartDate, t.Deadline, t.EstimatedHours, t.Status, t.Category, t.IsExternal, t.Description)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	t.ID = int(id)
	return nil
}

func (r *TaskRepo) Update(t *models.Task) error {
	var meetingID any
	if t.MeetingID != nil {
		meetingID = *t.MeetingID
	}
	_, err := r.db.Exec(`
		UPDATE tasks SET project_id=?, meeting_id=?, title=?, start_date=?, deadline=?,
			estimated_hours=?, status=?, category=?, is_external=?, description=?
		WHERE id=?`,
		t.ProjectID, meetingID, t.Title, t.StartDate, t.Deadline, t.EstimatedHours,
		t.Status, t.Category, t.IsExternal, t.Description, t.ID)
	return err
}

func (r *TaskRepo) UpdateStatus(id int, status string) error {
	_, err := r.db.Exec(`UPDATE tasks SET status = ? WHERE id = ?`, status, id)
	return err
}

func (r *TaskRepo) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM tasks WHERE id = ?`, id)
	return err
}

func (r *TaskRepo) ListByProject(projectID int) ([]models.Task, error) {
	return r.List(TaskFilter{ProjectID: projectID})
}

func (r *TaskRepo) ListByMeeting(meetingID int) ([]models.Task, error) {
	rows, err := r.db.Query(`
		SELECT t.id, t.title, t.deadline, t.status, t.category
		FROM tasks t WHERE t.meeting_id = ?
		ORDER BY t.deadline ASC`, meetingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Deadline, &t.Status, &t.Category); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (r *TaskRepo) ListOverdue() ([]models.Task, error) {
	return r.List(TaskFilter{Overdue: true})
}

// Dependencies
func (r *TaskRepo) AddDependency(taskID, dependsOnID int) error {
	_, err := r.db.Exec(`INSERT OR IGNORE INTO task_dependencies (task_id, depends_on_id) VALUES (?, ?)`, taskID, dependsOnID)
	return err
}

func (r *TaskRepo) RemoveDependency(taskID, dependsOnID int) error {
	_, err := r.db.Exec(`DELETE FROM task_dependencies WHERE task_id = ? AND depends_on_id = ?`, taskID, dependsOnID)
	return err
}

func (r *TaskRepo) GetDependencies(taskID int) ([]models.Task, error) {
	rows, err := r.db.Query(`
		SELECT t.id, t.title, t.status
		FROM task_dependencies td
		JOIN tasks t ON t.id = td.depends_on_id
		WHERE td.task_id = ?`, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Status); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}
