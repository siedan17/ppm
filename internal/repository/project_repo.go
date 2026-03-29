package repository

import (
	"database/sql"
	"fmt"

	"github.com/daniel/ppm/internal/database"
	"github.com/daniel/ppm/internal/models"
)

type ProjectRepo struct {
	db *database.DB
}

func NewProjectRepo(db *database.DB) *ProjectRepo {
	return &ProjectRepo{db: db}
}

func (r *ProjectRepo) List() ([]models.Project, error) {
	rows, err := r.db.Query(`
		SELECT p.id, p.name, p.priority, p.start_date, COALESCE(p.end_date,''), p.status,
			p.static_info, p.dynamic_info, p.created_at, p.updated_at
		FROM projects p
		ORDER BY p.priority ASC, p.name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var p models.Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Priority, &p.StartDate, &p.EndDate, &p.Status,
			&p.StaticInfo, &p.DynamicInfo, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, nil
}

func (r *ProjectRepo) ListActive() ([]models.Project, error) {
	rows, err := r.db.Query(`
		SELECT id, name, priority, start_date, COALESCE(end_date,''), status
		FROM projects WHERE status = 'active'
		ORDER BY priority ASC, name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var p models.Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Priority, &p.StartDate, &p.EndDate, &p.Status); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, nil
}

func (r *ProjectRepo) GetByID(id int) (*models.Project, error) {
	var p models.Project
	err := r.db.QueryRow(`
		SELECT id, name, priority, start_date, COALESCE(end_date,''), status,
			static_info, dynamic_info, created_at, updated_at
		FROM projects WHERE id = ?`, id).
		Scan(&p.ID, &p.Name, &p.Priority, &p.StartDate, &p.EndDate, &p.Status,
			&p.StaticInfo, &p.DynamicInfo, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("project not found")
	}
	if err != nil {
		return nil, err
	}

	// Load linked people
	pRows, err := r.db.Query(`
		SELECT pe.id, pe.name, pe.company, pe.role, COALESCE(pe.email,''), COALESCE(pe.phone,''),
			pp.role_in_project
		FROM project_people pp
		JOIN people pe ON pe.id = pp.person_id
		WHERE pp.project_id = ?
		ORDER BY pe.name`, id)
	if err != nil {
		return nil, err
	}
	defer pRows.Close()

	for pRows.Next() {
		var person models.Person
		var roleInProject string
		if err := pRows.Scan(&person.ID, &person.Name, &person.Company, &person.Role, &person.Email, &person.Phone, &roleInProject); err != nil {
			return nil, err
		}
		p.People = append(p.People, person)
	}

	return &p, nil
}

func (r *ProjectRepo) Create(p *models.Project) error {
	result, err := r.db.Exec(`INSERT INTO projects (name, priority, start_date, end_date, status, static_info, dynamic_info) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		p.Name, p.Priority, p.StartDate, nilIfEmpty(p.EndDate), p.Status, p.StaticInfo, p.DynamicInfo)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	p.ID = int(id)
	return nil
}

func (r *ProjectRepo) Update(p *models.Project) error {
	_, err := r.db.Exec(`UPDATE projects SET name=?, priority=?, start_date=?, end_date=?, status=?, static_info=?, dynamic_info=? WHERE id=?`,
		p.Name, p.Priority, p.StartDate, nilIfEmpty(p.EndDate), p.Status, p.StaticInfo, p.DynamicInfo, p.ID)
	return err
}

func (r *ProjectRepo) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM projects WHERE id = ?`, id)
	return err
}

func (r *ProjectRepo) LinkPerson(projectID, personID int, roleInProject string) error {
	_, err := r.db.Exec(`INSERT OR REPLACE INTO project_people (project_id, person_id, role_in_project) VALUES (?, ?, ?)`,
		projectID, personID, roleInProject)
	return err
}

func (r *ProjectRepo) UnlinkPerson(projectID, personID int) error {
	_, err := r.db.Exec(`DELETE FROM project_people WHERE project_id = ? AND person_id = ?`, projectID, personID)
	return err
}
