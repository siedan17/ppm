-- name: ListProjects :many
SELECT id, name, priority, start_date, COALESCE(end_date,'') AS end_date, status,
    static_info, dynamic_info, created_at, updated_at
FROM projects ORDER BY priority ASC, name ASC;

-- name: ListActiveProjects :many
SELECT id, name, priority, start_date, COALESCE(end_date,'') AS end_date, status
FROM projects WHERE status = 'active'
ORDER BY priority ASC, name ASC;

-- name: GetProjectByID :one
SELECT id, name, priority, start_date, COALESCE(end_date,'') AS end_date, status,
    static_info, dynamic_info, created_at, updated_at
FROM projects WHERE id = ?;

-- name: CreateProject :execlastid
INSERT INTO projects (name, priority, start_date, end_date, status, static_info, dynamic_info)
VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: UpdateProject :exec
UPDATE projects SET name=?, priority=?, start_date=?, end_date=?, status=?, static_info=?, dynamic_info=? WHERE id=?;

-- name: DeleteProject :exec
DELETE FROM projects WHERE id = ?;

-- name: LinkPersonToProject :exec
INSERT OR REPLACE INTO project_people (project_id, person_id, role_in_project) VALUES (?, ?, ?);

-- name: UnlinkPersonFromProject :exec
DELETE FROM project_people WHERE project_id = ? AND person_id = ?;

-- name: GetProjectPeople :many
SELECT pe.id, pe.name, pe.company, pe.role, COALESCE(pe.email,'') AS email, COALESCE(pe.phone,'') AS phone,
    pp.role_in_project
FROM project_people pp
JOIN people pe ON pe.id = pp.person_id
WHERE pp.project_id = ? ORDER BY pe.name;
