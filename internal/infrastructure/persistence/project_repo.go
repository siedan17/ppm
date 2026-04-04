package persistence

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/daniel/ppm/internal/domain"
	"github.com/daniel/ppm/internal/infrastructure/persistence/sqlcdb"
)

var _ domain.ProjectRepository = (*ProjectRepo)(nil)

type ProjectRepo struct {
	q *sqlcdb.Queries
}

func NewProjectRepo(q *sqlcdb.Queries) *ProjectRepo {
	return &ProjectRepo{q: q}
}

func (r *ProjectRepo) List() ([]domain.Project, error) {
	rows, err := r.q.ListProjects(context.Background())
	if err != nil {
		return nil, err
	}
	projects := make([]domain.Project, len(rows))
	for i, row := range rows {
		projects[i] = domain.Project{
			ID: int(row.ID), Name: row.Name, Priority: int(row.Priority),
			StartDate: row.StartDate, EndDate: row.EndDate, Status: row.Status,
			StaticInfo: row.StaticInfo, DynamicInfo: row.DynamicInfo,
			CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt,
		}
	}
	return projects, nil
}

func (r *ProjectRepo) ListActive() ([]domain.Project, error) {
	rows, err := r.q.ListActiveProjects(context.Background())
	if err != nil {
		return nil, err
	}
	projects := make([]domain.Project, len(rows))
	for i, row := range rows {
		projects[i] = domain.Project{
			ID: int(row.ID), Name: row.Name, Priority: int(row.Priority),
			StartDate: row.StartDate, EndDate: row.EndDate, Status: row.Status,
		}
	}
	return projects, nil
}

func (r *ProjectRepo) GetByID(id int) (*domain.Project, error) {
	row, err := r.q.GetProjectByID(context.Background(), int64(id))
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("project not found")
	}
	if err != nil {
		return nil, err
	}

	p := &domain.Project{
		ID: int(row.ID), Name: row.Name, Priority: int(row.Priority),
		StartDate: row.StartDate, EndDate: row.EndDate, Status: row.Status,
		StaticInfo: row.StaticInfo, DynamicInfo: row.DynamicInfo,
		CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt,
	}

	// Load linked people
	pRows, err := r.q.GetProjectPeople(context.Background(), int64(id))
	if err != nil {
		return nil, err
	}
	for _, pr := range pRows {
		p.People = append(p.People, domain.Person{
			ID: int(pr.ID), Name: pr.Name, Company: pr.Company,
			Role: pr.Role, Email: pr.Email, Phone: pr.Phone,
		})
	}

	return p, nil
}

func (r *ProjectRepo) Create(p *domain.Project) error {
	id, err := r.q.CreateProject(context.Background(), sqlcdb.CreateProjectParams{
		Name: p.Name, Priority: int64(p.Priority), StartDate: p.StartDate,
		EndDate: nullString(p.EndDate), Status: p.Status,
		StaticInfo: p.StaticInfo, DynamicInfo: p.DynamicInfo,
	})
	if err != nil {
		return err
	}
	p.ID = int(id)
	return nil
}

func (r *ProjectRepo) Update(p *domain.Project) error {
	return r.q.UpdateProject(context.Background(), sqlcdb.UpdateProjectParams{
		Name: p.Name, Priority: int64(p.Priority), StartDate: p.StartDate,
		EndDate: nullString(p.EndDate), Status: p.Status,
		StaticInfo: p.StaticInfo, DynamicInfo: p.DynamicInfo, ID: int64(p.ID),
	})
}

func (r *ProjectRepo) Delete(id int) error {
	return r.q.DeleteProject(context.Background(), int64(id))
}

func (r *ProjectRepo) LinkPerson(projectID, personID int, roleInProject string) error {
	return r.q.LinkPersonToProject(context.Background(), sqlcdb.LinkPersonToProjectParams{
		ProjectID: int64(projectID), PersonID: int64(personID), RoleInProject: roleInProject,
	})
}

func (r *ProjectRepo) UnlinkPerson(projectID, personID int) error {
	return r.q.UnlinkPersonFromProject(context.Background(), sqlcdb.UnlinkPersonFromProjectParams{
		ProjectID: int64(projectID), PersonID: int64(personID),
	})
}
