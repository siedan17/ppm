-- +goose Up
ALTER TABLE people ADD COLUMN person_type TEXT NOT NULL DEFAULT 'external' CHECK (person_type IN ('internal','external'));

-- +goose Down
ALTER TABLE people DROP COLUMN person_type;
