package handlers

import (
	"net/http"
	"strconv"

	"github.com/daniel/ppm/internal/models"
	"github.com/daniel/ppm/internal/render"
	"github.com/daniel/ppm/internal/repository"
)

type ProjectsHandler struct {
	repo        *repository.ProjectRepo
	personRepo  *repository.PersonRepo
	taskRepo    *repository.TaskRepo
	meetingRepo *repository.MeetingRepo
	render      *render.Renderer
}

func NewProjectsHandler(repo *repository.ProjectRepo, personRepo *repository.PersonRepo, taskRepo *repository.TaskRepo, meetingRepo *repository.MeetingRepo, r *render.Renderer) *ProjectsHandler {
	return &ProjectsHandler{repo: repo, personRepo: personRepo, taskRepo: taskRepo, meetingRepo: meetingRepo, render: r}
}

func (h *ProjectsHandler) List(w http.ResponseWriter, r *http.Request) {
	projects, err := h.repo.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.render.Page(w, http.StatusOK, "project_list.html", render.PageData{
		Title:   "Projects",
		Content: projects,
	})
}

func (h *ProjectsHandler) New(w http.ResponseWriter, r *http.Request) {
	h.render.Page(w, http.StatusOK, "project_form.html", render.PageData{
		Title: "New Project",
		Content: map[string]any{
			"Project":    &models.Project{Priority: 3, Status: "active"},
			"Statuses":   models.ProjectStatuses,
			"Priorities": models.Priorities,
		},
	})
}

func (h *ProjectsHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	priority, _ := strconv.Atoi(r.FormValue("priority"))
	p := &models.Project{
		Name:        r.FormValue("name"),
		Priority:    priority,
		StartDate:   r.FormValue("start_date"),
		EndDate:     r.FormValue("end_date"),
		Status:      r.FormValue("status"),
		StaticInfo:  r.FormValue("static_info"),
		DynamicInfo: r.FormValue("dynamic_info"),
	}

	errors := validateProject(p)
	if len(errors) > 0 {
		h.render.Page(w, http.StatusUnprocessableEntity, "project_form.html", render.PageData{
			Title: "New Project",
			Content: map[string]any{
				"Project":    p,
				"Statuses":   models.ProjectStatuses,
				"Priorities": models.Priorities,
				"Errors":     errors,
			},
			Flash: "Please fix the errors below.",
		})
		return
	}

	if err := h.repo.Create(p); err != nil {
		h.render.Page(w, http.StatusUnprocessableEntity, "project_form.html", render.PageData{
			Title: "New Project",
			Content: map[string]any{
				"Project":    p,
				"Statuses":   models.ProjectStatuses,
				"Priorities": models.Priorities,
			},
			Flash: "Error: " + err.Error(),
		})
		return
	}

	http.Redirect(w, r, "/projects/"+strconv.Itoa(p.ID), http.StatusSeeOther)
}

func (h *ProjectsHandler) Detail(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	project, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	tasks, _ := h.taskRepo.ListByProject(id)
	meetings, _ := h.meetingRepo.ListByProject(id)
	allPeople, _ := h.personRepo.List()

	h.render.Page(w, http.StatusOK, "project_detail.html", render.PageData{
		Title: project.Name,
		Content: map[string]any{
			"Project":   project,
			"Tasks":     tasks,
			"Meetings":  meetings,
			"AllPeople": allPeople,
		},
	})
}

func (h *ProjectsHandler) Edit(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	project, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	h.render.Page(w, http.StatusOK, "project_form.html", render.PageData{
		Title: "Edit " + project.Name,
		Content: map[string]any{
			"Project":    project,
			"Statuses":   models.ProjectStatuses,
			"Priorities": models.Priorities,
		},
	})
}

func (h *ProjectsHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	priority, _ := strconv.Atoi(r.FormValue("priority"))
	p := &models.Project{
		ID:          id,
		Name:        r.FormValue("name"),
		Priority:    priority,
		StartDate:   r.FormValue("start_date"),
		EndDate:     r.FormValue("end_date"),
		Status:      r.FormValue("status"),
		StaticInfo:  r.FormValue("static_info"),
		DynamicInfo: r.FormValue("dynamic_info"),
	}

	errors := validateProject(p)
	if len(errors) > 0 {
		h.render.Page(w, http.StatusUnprocessableEntity, "project_form.html", render.PageData{
			Title: "Edit " + p.Name,
			Content: map[string]any{
				"Project":    p,
				"Statuses":   models.ProjectStatuses,
				"Priorities": models.Priorities,
				"Errors":     errors,
			},
			Flash: "Please fix the errors below.",
		})
		return
	}

	if err := h.repo.Update(p); err != nil {
		h.render.Page(w, http.StatusUnprocessableEntity, "project_form.html", render.PageData{
			Title: "Edit " + p.Name,
			Content: map[string]any{
				"Project":    p,
				"Statuses":   models.ProjectStatuses,
				"Priorities": models.Priorities,
			},
			Flash: "Error: " + err.Error(),
		})
		return
	}

	http.Redirect(w, r, "/projects/"+strconv.Itoa(id), http.StatusSeeOther)
}

func (h *ProjectsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/projects")
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/projects", http.StatusSeeOther)
}

func (h *ProjectsHandler) LinkPerson(w http.ResponseWriter, r *http.Request) {
	projectID, _ := strconv.Atoi(r.PathValue("id"))
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	personID, _ := strconv.Atoi(r.FormValue("person_id"))
	roleInProject := r.FormValue("role_in_project")

	if err := h.repo.LinkPerson(projectID, personID, roleInProject); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		// Re-render the project detail to show updated people
		http.Redirect(w, r, "/projects/"+strconv.Itoa(projectID), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/projects/"+strconv.Itoa(projectID), http.StatusSeeOther)
}

func (h *ProjectsHandler) UnlinkPerson(w http.ResponseWriter, r *http.Request) {
	projectID, _ := strconv.Atoi(r.PathValue("id"))
	personID, _ := strconv.Atoi(r.PathValue("pid"))

	if err := h.repo.UnlinkPerson(projectID, personID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/projects/"+strconv.Itoa(projectID), http.StatusSeeOther)
}

func validateProject(p *models.Project) map[string]string {
	errors := make(map[string]string)
	if p.Name == "" {
		errors["name"] = "Name is required"
	}
	if p.StartDate == "" {
		errors["start_date"] = "Start date is required"
	}
	if p.Priority < 1 || p.Priority > 5 {
		errors["priority"] = "Priority must be between 1 and 5"
	}
	return errors
}
