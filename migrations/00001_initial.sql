-- +goose Up

CREATE TABLE projects (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    name          TEXT    NOT NULL UNIQUE,
    priority      INTEGER NOT NULL DEFAULT 3 CHECK (priority BETWEEN 1 AND 5),
    start_date    TEXT    NOT NULL,
    end_date      TEXT,
    status        TEXT    NOT NULL DEFAULT 'active'
                  CHECK (status IN ('active','on_hold','completed','archived')),
    static_info   TEXT    NOT NULL DEFAULT '',
    dynamic_info  TEXT    NOT NULL DEFAULT '',
    created_at    TEXT    NOT NULL DEFAULT (datetime('now')),
    updated_at    TEXT    NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE people (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    name       TEXT NOT NULL,
    company    TEXT NOT NULL DEFAULT '',
    role       TEXT NOT NULL DEFAULT '',
    email      TEXT,
    phone      TEXT,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE project_people (
    project_id      INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    person_id       INTEGER NOT NULL REFERENCES people(id) ON DELETE CASCADE,
    role_in_project TEXT NOT NULL DEFAULT '',
    PRIMARY KEY (project_id, person_id)
);

CREATE TABLE meetings (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    project_id   INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    date         TEXT    NOT NULL,
    meeting_type TEXT    NOT NULL CHECK (meeting_type IN ('internal','external')),
    title        TEXT    NOT NULL DEFAULT '',
    notes        TEXT    NOT NULL DEFAULT '',
    created_at   TEXT    NOT NULL DEFAULT (datetime('now')),
    updated_at   TEXT    NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE meeting_participants (
    meeting_id INTEGER NOT NULL REFERENCES meetings(id) ON DELETE CASCADE,
    person_id  INTEGER NOT NULL REFERENCES people(id) ON DELETE CASCADE,
    PRIMARY KEY (meeting_id, person_id)
);

CREATE TABLE tasks (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    project_id      INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    meeting_id      INTEGER REFERENCES meetings(id) ON DELETE SET NULL,
    title           TEXT    NOT NULL,
    start_date      TEXT    NOT NULL,
    deadline        TEXT    NOT NULL,
    estimated_hours REAL    NOT NULL CHECK (estimated_hours > 0),
    status          TEXT    NOT NULL DEFAULT 'todo'
                    CHECK (status IN ('todo','in_progress','blocked','done','cancelled')),
    category        TEXT    NOT NULL DEFAULT 'other'
                    CHECK (category IN ('programming','data_engineering','specification',
                                        'design','communication','other')),
    is_external     INTEGER NOT NULL DEFAULT 0 CHECK (is_external IN (0,1)),
    description     TEXT    NOT NULL DEFAULT '',
    created_at      TEXT    NOT NULL DEFAULT (datetime('now')),
    updated_at      TEXT    NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE task_dependencies (
    task_id       INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    depends_on_id INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    PRIMARY KEY (task_id, depends_on_id),
    CHECK (task_id != depends_on_id)
);

CREATE INDEX idx_projects_priority ON projects(priority);
CREATE INDEX idx_tasks_project_id  ON tasks(project_id);
CREATE INDEX idx_tasks_status      ON tasks(status);
CREATE INDEX idx_tasks_deadline    ON tasks(deadline);
CREATE INDEX idx_meetings_project  ON meetings(project_id);
CREATE INDEX idx_meetings_date     ON meetings(date);

-- +goose StatementBegin
CREATE TRIGGER trg_projects_updated AFTER UPDATE ON projects
BEGIN UPDATE projects SET updated_at = datetime('now') WHERE id = NEW.id; END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_people_updated AFTER UPDATE ON people
BEGIN UPDATE people SET updated_at = datetime('now') WHERE id = NEW.id; END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_meetings_updated AFTER UPDATE ON meetings
BEGIN UPDATE meetings SET updated_at = datetime('now') WHERE id = NEW.id; END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_tasks_updated AFTER UPDATE ON tasks
BEGIN UPDATE tasks SET updated_at = datetime('now') WHERE id = NEW.id; END;
-- +goose StatementEnd

-- +goose Down

DROP TRIGGER IF EXISTS trg_tasks_updated;
DROP TRIGGER IF EXISTS trg_meetings_updated;
DROP TRIGGER IF EXISTS trg_people_updated;
DROP TRIGGER IF EXISTS trg_projects_updated;
DROP INDEX IF EXISTS idx_meetings_date;
DROP INDEX IF EXISTS idx_meetings_project;
DROP INDEX IF EXISTS idx_tasks_deadline;
DROP INDEX IF EXISTS idx_tasks_status;
DROP INDEX IF EXISTS idx_tasks_project_id;
DROP INDEX IF EXISTS idx_projects_priority;
DROP TABLE IF EXISTS task_dependencies;
DROP TABLE IF EXISTS meeting_participants;
DROP TABLE IF EXISTS tasks;
DROP TABLE IF EXISTS meetings;
DROP TABLE IF EXISTS project_people;
DROP TABLE IF EXISTS people;
DROP TABLE IF EXISTS projects;
