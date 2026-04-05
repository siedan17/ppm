-- name: ListPeople :many
SELECT id, name, company, role, COALESCE(email,'') AS email, COALESCE(phone,'') AS phone, person_type, created_at, updated_at
FROM people ORDER BY name ASC;

-- name: GetPersonByID :one
SELECT id, name, company, role, COALESCE(email,'') AS email, COALESCE(phone,'') AS phone, person_type, created_at, updated_at
FROM people WHERE id = ?;

-- name: CreatePerson :execlastid
INSERT INTO people (name, company, role, email, phone, person_type) VALUES (?, ?, ?, ?, ?, ?);

-- name: UpdatePerson :exec
UPDATE people SET name=?, company=?, role=?, email=?, phone=?, person_type=? WHERE id=?;

-- name: DeletePerson :exec
DELETE FROM people WHERE id = ?;
