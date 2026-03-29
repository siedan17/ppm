package repository

import (
	"database/sql"
	"fmt"

	"github.com/daniel/ppm/internal/database"
	"github.com/daniel/ppm/internal/models"
)

type MeetingRepo struct {
	db *database.DB
}

func NewMeetingRepo(db *database.DB) *MeetingRepo {
	return &MeetingRepo{db: db}
}

type MeetingFilter struct {
	ProjectID int
	DateFrom  string
	DateTo    string
}

func (r *MeetingRepo) List(f MeetingFilter) ([]models.Meeting, error) {
	query := `
		SELECT m.id, m.project_id, m.date, m.meeting_type, m.title, m.notes,
			m.created_at, m.updated_at, p.name as project_name
		FROM meetings m
		JOIN projects p ON p.id = m.project_id
		WHERE 1=1`

	var args []any

	if f.ProjectID > 0 {
		query += " AND m.project_id = ?"
		args = append(args, f.ProjectID)
	}
	if f.DateFrom != "" {
		query += " AND m.date >= ?"
		args = append(args, f.DateFrom)
	}
	if f.DateTo != "" {
		query += " AND m.date <= ?"
		args = append(args, f.DateTo)
	}

	query += " ORDER BY m.date DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var meetings []models.Meeting
	for rows.Next() {
		var m models.Meeting
		if err := rows.Scan(&m.ID, &m.ProjectID, &m.Date, &m.MeetingType, &m.Title, &m.Notes,
			&m.CreatedAt, &m.UpdatedAt, &m.ProjectName); err != nil {
			return nil, err
		}
		meetings = append(meetings, m)
	}
	return meetings, nil
}

func (r *MeetingRepo) GetByID(id int) (*models.Meeting, error) {
	var m models.Meeting
	err := r.db.QueryRow(`
		SELECT m.id, m.project_id, m.date, m.meeting_type, m.title, m.notes,
			m.created_at, m.updated_at, p.name as project_name
		FROM meetings m
		JOIN projects p ON p.id = m.project_id
		WHERE m.id = ?`, id).
		Scan(&m.ID, &m.ProjectID, &m.Date, &m.MeetingType, &m.Title, &m.Notes,
			&m.CreatedAt, &m.UpdatedAt, &m.ProjectName)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("meeting not found")
	}
	if err != nil {
		return nil, err
	}

	// Load participants
	pRows, err := r.db.Query(`
		SELECT pe.id, pe.name, pe.company, pe.role, COALESCE(pe.email,''), COALESCE(pe.phone,'')
		FROM meeting_participants mp
		JOIN people pe ON pe.id = mp.person_id
		WHERE mp.meeting_id = ?
		ORDER BY pe.name`, id)
	if err != nil {
		return nil, err
	}
	defer pRows.Close()

	for pRows.Next() {
		var p models.Person
		if err := pRows.Scan(&p.ID, &p.Name, &p.Company, &p.Role, &p.Email, &p.Phone); err != nil {
			return nil, err
		}
		m.Participants = append(m.Participants, p)
	}

	return &m, nil
}

func (r *MeetingRepo) Create(m *models.Meeting) error {
	result, err := r.db.Exec(`INSERT INTO meetings (project_id, date, meeting_type, title, notes) VALUES (?, ?, ?, ?, ?)`,
		m.ProjectID, m.Date, m.MeetingType, m.Title, m.Notes)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	m.ID = int(id)
	return nil
}

func (r *MeetingRepo) Update(m *models.Meeting) error {
	_, err := r.db.Exec(`UPDATE meetings SET project_id=?, date=?, meeting_type=?, title=?, notes=? WHERE id=?`,
		m.ProjectID, m.Date, m.MeetingType, m.Title, m.Notes, m.ID)
	return err
}

func (r *MeetingRepo) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM meetings WHERE id = ?`, id)
	return err
}

func (r *MeetingRepo) AddParticipant(meetingID, personID int) error {
	_, err := r.db.Exec(`INSERT OR IGNORE INTO meeting_participants (meeting_id, person_id) VALUES (?, ?)`, meetingID, personID)
	return err
}

func (r *MeetingRepo) RemoveParticipant(meetingID, personID int) error {
	_, err := r.db.Exec(`DELETE FROM meeting_participants WHERE meeting_id = ? AND person_id = ?`, meetingID, personID)
	return err
}

func (r *MeetingRepo) ListByProject(projectID int) ([]models.Meeting, error) {
	return r.List(MeetingFilter{ProjectID: projectID})
}
