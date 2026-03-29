package repository

import (
	"database/sql"
	"fmt"

	"github.com/daniel/ppm/internal/database"
	"github.com/daniel/ppm/internal/models"
)

type PersonRepo struct {
	db *database.DB
}

func NewPersonRepo(db *database.DB) *PersonRepo {
	return &PersonRepo{db: db}
}

func (r *PersonRepo) List() ([]models.Person, error) {
	rows, err := r.db.Query(`SELECT id, name, company, role, COALESCE(email,''), COALESCE(phone,''), created_at, updated_at FROM people ORDER BY name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var people []models.Person
	for rows.Next() {
		var p models.Person
		if err := rows.Scan(&p.ID, &p.Name, &p.Company, &p.Role, &p.Email, &p.Phone, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		people = append(people, p)
	}
	return people, nil
}

func (r *PersonRepo) GetByID(id int) (*models.Person, error) {
	var p models.Person
	err := r.db.QueryRow(`SELECT id, name, company, role, COALESCE(email,''), COALESCE(phone,''), created_at, updated_at FROM people WHERE id = ?`, id).
		Scan(&p.ID, &p.Name, &p.Company, &p.Role, &p.Email, &p.Phone, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("person not found")
	}
	return &p, err
}

func (r *PersonRepo) Create(p *models.Person) error {
	result, err := r.db.Exec(`INSERT INTO people (name, company, role, email, phone) VALUES (?, ?, ?, ?, ?)`,
		p.Name, p.Company, p.Role, nilIfEmpty(p.Email), nilIfEmpty(p.Phone))
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	p.ID = int(id)
	return nil
}

func (r *PersonRepo) Update(p *models.Person) error {
	_, err := r.db.Exec(`UPDATE people SET name=?, company=?, role=?, email=?, phone=? WHERE id=?`,
		p.Name, p.Company, p.Role, nilIfEmpty(p.Email), nilIfEmpty(p.Phone), p.ID)
	return err
}

func (r *PersonRepo) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM people WHERE id = ?`, id)
	return err
}

func nilIfEmpty(s string) any {
	if s == "" {
		return nil
	}
	return s
}
