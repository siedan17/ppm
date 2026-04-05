-- +goose Up
ALTER TABLE projects ADD COLUMN general_info TEXT NOT NULL DEFAULT '';
ALTER TABLE projects ADD COLUMN timeline TEXT NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE projects DROP COLUMN general_info;
ALTER TABLE projects DROP COLUMN timeline;
