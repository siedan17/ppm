-- name: ListMeetings :many
SELECT m.id, m.project_id, m.date, m.meeting_type, m.title, m.notes,
    m.created_at, m.updated_at, p.name AS project_name
FROM meetings m
JOIN projects p ON p.id = m.project_id
WHERE (sqlc.narg('project_id') IS NULL OR m.project_id = sqlc.narg('project_id'))
  AND (sqlc.narg('date_from') IS NULL OR m.date >= sqlc.narg('date_from'))
  AND (sqlc.narg('date_to') IS NULL OR m.date <= sqlc.narg('date_to'))
ORDER BY m.date DESC;

-- name: GetMeetingByID :one
SELECT m.id, m.project_id, m.date, m.meeting_type, m.title, m.notes,
    m.created_at, m.updated_at, p.name AS project_name
FROM meetings m
JOIN projects p ON p.id = m.project_id
WHERE m.id = ?;

-- name: CreateMeeting :execlastid
INSERT INTO meetings (project_id, date, meeting_type, title, notes) VALUES (?, ?, ?, ?, ?);

-- name: UpdateMeeting :exec
UPDATE meetings SET project_id=?, date=?, meeting_type=?, title=?, notes=? WHERE id=?;

-- name: DeleteMeeting :exec
DELETE FROM meetings WHERE id = ?;

-- name: AddMeetingParticipant :exec
INSERT OR IGNORE INTO meeting_participants (meeting_id, person_id) VALUES (?, ?);

-- name: RemoveMeetingParticipant :exec
DELETE FROM meeting_participants WHERE meeting_id = ? AND person_id = ?;

-- name: GetMeetingParticipants :many
SELECT pe.id, pe.name, pe.company, pe.role, COALESCE(pe.email,'') AS email, COALESCE(pe.phone,'') AS phone
FROM meeting_participants mp
JOIN people pe ON pe.id = mp.person_id
WHERE mp.meeting_id = ?
ORDER BY pe.name;

-- name: ListMeetingsByProject :many
SELECT m.id, m.project_id, m.date, m.meeting_type, m.title, m.notes,
    m.created_at, m.updated_at, p.name AS project_name
FROM meetings m
JOIN projects p ON p.id = m.project_id
WHERE m.project_id = ?
ORDER BY m.date DESC;
