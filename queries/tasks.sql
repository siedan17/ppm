-- name: ListTasks :many
SELECT t.id, t.project_id, t.meeting_id, t.title, t.start_date, t.deadline,
    t.estimated_hours, t.status, t.category, t.is_external, t.description,
    t.created_at, t.updated_at, p.name AS project_name,
    COALESCE(m.title, '') AS meeting_title,
    CASE WHEN t.deadline < date('now') AND t.status NOT IN ('done','cancelled') THEN 1 ELSE 0 END AS is_overdue
FROM tasks t
JOIN projects p ON p.id = t.project_id
LEFT JOIN meetings m ON m.id = t.meeting_id
WHERE (sqlc.narg('project_id') IS NULL OR t.project_id = sqlc.narg('project_id'))
  AND (sqlc.narg('status') IS NULL OR t.status = sqlc.narg('status'))
  AND (sqlc.narg('category') IS NULL OR t.category = sqlc.narg('category'))
  AND (sqlc.narg('overdue') IS NULL OR (t.deadline < date('now') AND t.status NOT IN ('done','cancelled')))
ORDER BY t.deadline ASC, t.status ASC;

-- name: GetTaskByID :one
SELECT t.id, t.project_id, t.meeting_id, t.title, t.start_date, t.deadline,
    t.estimated_hours, t.status, t.category, t.is_external, t.description,
    t.created_at, t.updated_at, p.name AS project_name,
    COALESCE(m.title, '') AS meeting_title
FROM tasks t
JOIN projects p ON p.id = t.project_id
LEFT JOIN meetings m ON m.id = t.meeting_id
WHERE t.id = ?;

-- name: CreateTask :execlastid
INSERT INTO tasks (project_id, meeting_id, title, start_date, deadline, estimated_hours, status, category, is_external, description)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: UpdateTask :exec
UPDATE tasks SET project_id=?, meeting_id=?, title=?, start_date=?, deadline=?,
    estimated_hours=?, status=?, category=?, is_external=?, description=?
WHERE id=?;

-- name: UpdateTaskStatus :exec
UPDATE tasks SET status = ? WHERE id = ?;

-- name: DeleteTask :exec
DELETE FROM tasks WHERE id = ?;

-- name: ListTasksByMeeting :many
SELECT t.id, t.title, t.deadline, t.status, t.category
FROM tasks t WHERE t.meeting_id = ?
ORDER BY t.deadline ASC;

-- name: AddTaskDependency :exec
INSERT OR IGNORE INTO task_dependencies (task_id, depends_on_id) VALUES (?, ?);

-- name: RemoveTaskDependency :exec
DELETE FROM task_dependencies WHERE task_id = ? AND depends_on_id = ?;

-- name: GetTaskDependencies :many
SELECT t.id, t.title, t.status
FROM task_dependencies td
JOIN tasks t ON t.id = td.depends_on_id
WHERE td.task_id = ?;
