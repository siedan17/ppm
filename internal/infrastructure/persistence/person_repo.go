package persistence

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/daniel/ppm/internal/domain"
	"github.com/daniel/ppm/internal/infrastructure/persistence/sqlcdb"
)

var _ domain.PersonRepository = (*PersonRepo)(nil)

type PersonRepo struct {
	q *sqlcdb.Queries
}

func NewPersonRepo(q *sqlcdb.Queries) *PersonRepo {
	return &PersonRepo{q: q}
}

func (r *PersonRepo) List() ([]domain.Person, error) {
	rows, err := r.q.ListPeople(context.Background())
	if err != nil {
		return nil, err
	}
	people := make([]domain.Person, len(rows))
	for i, row := range rows {
		people[i] = domain.Person{
			ID: int(row.ID), Name: row.Name, Company: row.Company, Role: row.Role,
			Email: row.Email, Phone: row.Phone, CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt,
		}
	}
	return people, nil
}

func (r *PersonRepo) GetByID(id int) (*domain.Person, error) {
	row, err := r.q.GetPersonByID(context.Background(), int64(id))
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("person not found")
	}
	if err != nil {
		return nil, err
	}
	return &domain.Person{
		ID: int(row.ID), Name: row.Name, Company: row.Company, Role: row.Role,
		Email: row.Email, Phone: row.Phone, CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt,
	}, nil
}

func (r *PersonRepo) Create(p *domain.Person) error {
	id, err := r.q.CreatePerson(context.Background(), sqlcdb.CreatePersonParams{
		Name: p.Name, Company: p.Company, Role: p.Role,
		Email: nullString(p.Email), Phone: nullString(p.Phone),
	})
	if err != nil {
		return err
	}
	p.ID = int(id)
	return nil
}

func (r *PersonRepo) Update(p *domain.Person) error {
	return r.q.UpdatePerson(context.Background(), sqlcdb.UpdatePersonParams{
		Name: p.Name, Company: p.Company, Role: p.Role,
		Email: nullString(p.Email), Phone: nullString(p.Phone), ID: int64(p.ID),
	})
}

func (r *PersonRepo) Delete(id int) error {
	return r.q.DeletePerson(context.Background(), int64(id))
}

func nullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}
