-- name: GetDashboardProjects :many
SELECT p.id, p.name, p.priority, p.status,
    COALESCE((SELECT COUNT(*) FROM tasks t WHERE t.project_id = p.id AND t.status NOT IN ('done','cancelled')), 0) AS task_count,
    COALESCE((SELECT COUNT(*) FROM tasks t WHERE t.project_id = p.id AND t.deadline < date('now') AND t.status NOT IN ('done','cancelled')), 0) AS overdue_count
FROM projects p
WHERE p.status = 'active'
ORDER BY p.priority ASC, p.name ASC;

-- name: GetActiveTasks :many
SELECT t.id, t.title, t.deadline, t.status, t.category, t.project_id, p.name AS project_name,
    CASE WHEN t.deadline < date('now') AND t.status NOT IN ('done','cancelled') THEN 1 ELSE 0 END AS is_overdue
FROM tasks t
JOIN projects p ON p.id = t.project_id
WHERE t.status NOT IN ('done','cancelled')
ORDER BY t.deadline ASC;

-- name: GetUpcomingMeetings :many
SELECT m.id, m.title, m.date, m.project_id, p.name AS project_name
FROM meetings m
JOIN projects p ON p.id = m.project_id
WHERE m.date BETWEEN date('now') AND date('now', '+7 days')
ORDER BY m.date ASC;
