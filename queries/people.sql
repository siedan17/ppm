-- name: ListPeople :many
SELECT id, name, company, role, COALESCE(email,'') AS email, COALESCE(phone,'') AS phone, created_at, updated_at
FROM people ORDER BY name ASC;

-- name: GetPersonByID :one
SELECT id, name, company, role, COALESCE(email,'') AS email, COALESCE(phone,'') AS phone, created_at, updated_at
FROM people WHERE id = ?;

-- name: CreatePerson :execlastid
INSERT INTO people (name, company, role, email, phone) VALUES (?, ?, ?, ?, ?);

-- name: UpdatePerson :exec
UPDATE people SET name=?, company=?, role=?, email=?, phone=? WHERE id=?;

-- name: DeletePerson :exec
DELETE FROM people WHERE id = ?;
