package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	ppm "github.com/daniel/ppm"
	"github.com/daniel/ppm/internal/application"
	deliveryhttp "github.com/daniel/ppm/internal/delivery/http"
	"github.com/daniel/ppm/internal/delivery/render"
	"github.com/daniel/ppm/internal/infrastructure/database"
	"github.com/daniel/ppm/internal/infrastructure/persistence"
	"github.com/daniel/ppm/internal/infrastructure/persistence/sqlcdb"
)

func main() {
	port := flag.Int("port", 8080, "HTTP port")
	dbPath := flag.String("db", "ppm.db", "SQLite database path")
	flag.Parse()

	// Open database with embedded migrations
	migrationsSub, err := fs.Sub(ppm.MigrationsFS, "migrations")
	if err != nil {
		log.Fatalf("Failed to get migrations sub-fs: %v", err)
	}
	db, err := database.Open(*dbPath, migrationsSub)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Initialize renderer
	renderer, err := render.New(ppm.TemplatesFS)
	if err != nil {
		log.Fatalf("Failed to initialize templates: %v", err)
	}

	// Infrastructure: sqlc queries and repository implementations
	queries := sqlcdb.New(db)
	personRepo := persistence.NewPersonRepo(queries)
	projectRepo := persistence.NewProjectRepo(queries)
	taskRepo := persistence.NewTaskRepo(queries)
	meetingRepo := persistence.NewMeetingRepo(queries)
	dashboardRepo := persistence.NewDashboardRepo(queries)

	// Application services
	personSvc := application.NewPersonService(personRepo)
	projectSvc := application.NewProjectService(projectRepo, taskRepo, meetingRepo, personRepo)
	taskSvc := application.NewTaskService(taskRepo, projectRepo)
	meetingSvc := application.NewMeetingService(meetingRepo, projectRepo, personRepo, taskRepo, renderer)
	dashboardSvc := application.NewDashboardService(dashboardRepo, taskRepo)
	exportSvc := application.NewExportService(meetingRepo, projectRepo)

	// Delivery: HTTP handlers
	dashboardH := deliveryhttp.NewDashboardHandler(dashboardSvc, renderer)
	peopleH := deliveryhttp.NewPeopleHandler(personSvc, renderer)
	projectsH := deliveryhttp.NewProjectsHandler(projectSvc, renderer)
	tasksH := deliveryhttp.NewTasksHandler(taskSvc, renderer)
	meetingsH := deliveryhttp.NewMeetingsHandler(meetingSvc, renderer)
	exportH := deliveryhttp.NewExportHandler(exportSvc)

	// Static files
	staticSub, err := fs.Sub(ppm.StaticFS, "static")
	if err != nil {
		log.Fatalf("Failed to get static sub-fs: %v", err)
	}

	mux := http.NewServeMux()

	// Static
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticSub))))

	// Dashboard
	mux.HandleFunc("GET /{$}", dashboardH.Index)
	mux.HandleFunc("GET /partials/dashboard-tasks", dashboardH.ActiveTasks)

	// People
	mux.HandleFunc("GET /people", peopleH.List)
	mux.HandleFunc("GET /people/new", peopleH.New)
	mux.HandleFunc("POST /people", peopleH.Create)
	mux.HandleFunc("GET /people/{id}/edit", peopleH.Edit)
	mux.HandleFunc("PUT /people/{id}", peopleH.Update)
	mux.HandleFunc("DELETE /people/{id}", peopleH.Delete)

	// Projects
	mux.HandleFunc("GET /projects", projectsH.List)
	mux.HandleFunc("GET /projects/new", projectsH.New)
	mux.HandleFunc("POST /projects", projectsH.Create)
	mux.HandleFunc("GET /projects/{id}", projectsH.Detail)
	mux.HandleFunc("GET /projects/{id}/edit", projectsH.Edit)
	mux.HandleFunc("PUT /projects/{id}", projectsH.Update)
	mux.HandleFunc("DELETE /projects/{id}", projectsH.Delete)
	mux.HandleFunc("POST /projects/{id}/people", projectsH.LinkPerson)
	mux.HandleFunc("DELETE /projects/{id}/people/{pid}", projectsH.UnlinkPerson)
	mux.HandleFunc("GET /projects/{id}/export/meetings", exportH.ExportProjectMeetings)

	// Tasks
	mux.HandleFunc("GET /tasks", tasksH.List)
	mux.HandleFunc("GET /tasks/new", tasksH.New)
	mux.HandleFunc("POST /tasks", tasksH.Create)
	mux.HandleFunc("GET /tasks/{id}/edit", tasksH.Edit)
	mux.HandleFunc("PUT /tasks/{id}", tasksH.Update)
	mux.HandleFunc("DELETE /tasks/{id}", tasksH.Delete)
	mux.HandleFunc("PATCH /tasks/{id}/status", tasksH.UpdateStatus)
	mux.HandleFunc("GET /partials/task-list", tasksH.TaskListPartial)

	// Meetings
	mux.HandleFunc("GET /meetings", meetingsH.List)
	mux.HandleFunc("GET /meetings/new", meetingsH.New)
	mux.HandleFunc("POST /meetings", meetingsH.Create)
	mux.HandleFunc("GET /meetings/{id}", meetingsH.Detail)
	mux.HandleFunc("GET /meetings/{id}/edit", meetingsH.Edit)
	mux.HandleFunc("PUT /meetings/{id}", meetingsH.Update)
	mux.HandleFunc("DELETE /meetings/{id}", meetingsH.Delete)
	mux.HandleFunc("GET /meetings/{id}/create-task", meetingsH.CreateTaskFromMeeting)
	mux.HandleFunc("POST /meetings/{id}/tasks", tasksH.Create)
	mux.HandleFunc("POST /meetings/{id}/participants", meetingsH.AddParticipant)
	mux.HandleFunc("DELETE /meetings/{id}/participants/{pid}", meetingsH.RemoveParticipant)
	mux.HandleFunc("GET /meetings/{id}/export", exportH.ExportMeeting)

	// Apply middleware
	handler := deliveryhttp.LoggingMiddleware(deliveryhttp.MethodOverride(mux))

	addr := fmt.Sprintf(":%d", *port)
	fmt.Fprintf(os.Stderr, "PPM running at http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, handler))
}
