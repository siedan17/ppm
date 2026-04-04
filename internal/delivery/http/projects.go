package http

import (
	"net/http"
	"strconv"

	"github.com/daniel/ppm/internal/application"
	"github.com/daniel/ppm/internal/delivery/render"
	"github.com/daniel/ppm/internal/domain"
)

type ProjectsHandler struct {
	svc    *application.ProjectService
	render *render.Renderer
}

func NewProjectsHandler(svc *application.ProjectService, r *render.Renderer) *ProjectsHandler {
	return &ProjectsHandler{svc: svc, render: r}
}

func (h *ProjectsHandler) List(w http.ResponseWriter, r *http.Request) {
	projects, err := h.svc.List()
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
			"Project":    &domain.Project{Priority: 3, Status: "active"},
			"Statuses":   domain.ProjectStatuses,
			"Priorities": domain.Priorities,
		},
	})
}

func (h *ProjectsHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	priority, _ := strconv.Atoi(r.FormValue("priority"))
	p := &domain.Project{
		Name:        r.FormValue("name"),
		Priority:    priority,
		StartDate:   r.FormValue("start_date"),
		EndDate:     r.FormValue("end_date"),
		Status:      r.FormValue("status"),
		StaticInfo:  r.FormValue("static_info"),
		DynamicInfo: r.FormValue("dynamic_info"),
	}

	errors := domain.ValidateProject(p)
	if len(errors) > 0 {
		h.render.Page(w, http.StatusUnprocessableEntity, "project_form.html", render.PageData{
			Title: "New Project",
			Content: map[string]any{
				"Project":    p,
				"Statuses":   domain.ProjectStatuses,
				"Priorities": domain.Priorities,
				"Errors":     errors,
			},
			Flash: "Please fix the errors below.",
		})
		return
	}

	if err := h.svc.Create(p); err != nil {
		h.render.Page(w, http.StatusUnprocessableEntity, "project_form.html", render.PageData{
			Title: "New Project",
			Content: map[string]any{
				"Project":    p,
				"Statuses":   domain.ProjectStatuses,
				"Priorities": domain.Priorities,
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

	detail, err := h.svc.GetDetail(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	h.render.Page(w, http.StatusOK, "project_detail.html", render.PageData{
		Title: detail.Project.Name,
		Content: map[string]any{
			"Project":   detail.Project,
			"Tasks":     detail.Tasks,
			"Meetings":  detail.Meetings,
			"AllPeople": detail.AllPeople,
		},
	})
}

func (h *ProjectsHandler) Edit(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	project, err := h.svc.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	h.render.Page(w, http.StatusOK, "project_form.html", render.PageData{
		Title: "Edit " + project.Name,
		Content: map[string]any{
			"Project":    project,
			"Statuses":   domain.ProjectStatuses,
			"Priorities": domain.Priorities,
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
	p := &domain.Project{
		ID:          id,
		Name:        r.FormValue("name"),
		Priority:    priority,
		StartDate:   r.FormValue("start_date"),
		EndDate:     r.FormValue("end_date"),
		Status:      r.FormValue("status"),
		StaticInfo:  r.FormValue("static_info"),
		DynamicInfo: r.FormValue("dynamic_info"),
	}

	errors := domain.ValidateProject(p)
	if len(errors) > 0 {
		h.render.Page(w, http.StatusUnprocessableEntity, "project_form.html", render.PageData{
			Title: "Edit " + p.Name,
			Content: map[string]any{
				"Project":    p,
				"Statuses":   domain.ProjectStatuses,
				"Priorities": domain.Priorities,
				"Errors":     errors,
			},
			Flash: "Please fix the errors below.",
		})
		return
	}

	if err := h.svc.Update(p); err != nil {
		h.render.Page(w, http.StatusUnprocessableEntity, "project_form.html", render.PageData{
			Title: "Edit " + p.Name,
			Content: map[string]any{
				"Project":    p,
				"Statuses":   domain.ProjectStatuses,
				"Priorities": domain.Priorities,
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

	if err := h.svc.Delete(id); err != nil {
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

	if err := h.svc.LinkPerson(projectID, personID, roleInProject); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/projects/"+strconv.Itoa(projectID), http.StatusSeeOther)
}

func (h *ProjectsHandler) UnlinkPerson(w http.ResponseWriter, r *http.Request) {
	projectID, _ := strconv.Atoi(r.PathValue("id"))
	personID, _ := strconv.Atoi(r.PathValue("pid"))

	if err := h.svc.UnlinkPerson(projectID, personID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/projects/"+strconv.Itoa(projectID), http.StatusSeeOther)
}
