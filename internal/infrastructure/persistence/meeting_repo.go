package persistence

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/daniel/ppm/internal/domain"
	"github.com/daniel/ppm/internal/infrastructure/persistence/sqlcdb"
)

var _ domain.MeetingRepository = (*MeetingRepo)(nil)

type MeetingRepo struct {
	q *sqlcdb.Queries
}

func NewMeetingRepo(q *sqlcdb.Queries) *MeetingRepo {
	return &MeetingRepo{q: q}
}

func (r *MeetingRepo) List(f domain.MeetingFilter) ([]domain.Meeting, error) {
	params := sqlcdb.ListMeetingsParams{}

	if f.ProjectID > 0 {
		params.ProjectID = int64(f.ProjectID)
	}
	if f.DateFrom != "" {
		params.DateFrom = f.DateFrom
	}
	if f.DateTo != "" {
		params.DateTo = f.DateTo
	}

	rows, err := r.q.ListMeetings(context.Background(), params)
	if err != nil {
		return nil, err
	}

	meetings := make([]domain.Meeting, len(rows))
	for i, row := range rows {
		meetings[i] = domain.Meeting{
			ID: int(row.ID), ProjectID: int(row.ProjectID), Date: row.Date,
			MeetingType: row.MeetingType, Title: row.Title, Notes: row.Notes,
			CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt,
			ProjectName: row.ProjectName,
		}
	}
	return meetings, nil
}

func (r *MeetingRepo) GetByID(id int) (*domain.Meeting, error) {
	row, err := r.q.GetMeetingByID(context.Background(), int64(id))
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("meeting not found")
	}
	if err != nil {
		return nil, err
	}

	m := &domain.Meeting{
		ID: int(row.ID), ProjectID: int(row.ProjectID), Date: row.Date,
		MeetingType: row.MeetingType, Title: row.Title, Notes: row.Notes,
		CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt,
		ProjectName: row.ProjectName,
	}

	// Load participants
	pRows, err := r.q.GetMeetingParticipants(context.Background(), int64(id))
	if err != nil {
		return nil, err
	}
	for _, pr := range pRows {
		m.Participants = append(m.Participants, domain.Person{
			ID: int(pr.ID), Name: pr.Name, Company: pr.Company,
			Role: pr.Role, Email: pr.Email, Phone: pr.Phone,
		})
	}

	return m, nil
}

func (r *MeetingRepo) Create(m *domain.Meeting) error {
	id, err := r.q.CreateMeeting(context.Background(), sqlcdb.CreateMeetingParams{
		ProjectID: int64(m.ProjectID), Date: m.Date,
		MeetingType: m.MeetingType, Title: m.Title, Notes: m.Notes,
	})
	if err != nil {
		return err
	}
	m.ID = int(id)
	return nil
}

func (r *MeetingRepo) Update(m *domain.Meeting) error {
	return r.q.UpdateMeeting(context.Background(), sqlcdb.UpdateMeetingParams{
		ProjectID: int64(m.ProjectID), Date: m.Date,
		MeetingType: m.MeetingType, Title: m.Title, Notes: m.Notes, ID: int64(m.ID),
	})
}

func (r *MeetingRepo) Delete(id int) error {
	return r.q.DeleteMeeting(context.Background(), int64(id))
}

func (r *MeetingRepo) AddParticipant(meetingID, personID int) error {
	return r.q.AddMeetingParticipant(context.Background(), sqlcdb.AddMeetingParticipantParams{
		MeetingID: int64(meetingID), PersonID: int64(personID),
	})
}

func (r *MeetingRepo) RemoveParticipant(meetingID, personID int) error {
	return r.q.RemoveMeetingParticipant(context.Background(), sqlcdb.RemoveMeetingParticipantParams{
		MeetingID: int64(meetingID), PersonID: int64(personID),
	})
}

func (r *MeetingRepo) ListByProject(projectID int) ([]domain.Meeting, error) {
	return r.List(domain.MeetingFilter{ProjectID: projectID})
}
