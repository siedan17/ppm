package http

import (
	"net/http"
	"strconv"

	"github.com/daniel/ppm/internal/application"
	"github.com/daniel/ppm/internal/delivery/render"
	"github.com/daniel/ppm/internal/domain"
)

type TasksHandler struct {
	svc    *application.TaskService
	render *render.Renderer
}

func NewTasksHandler(svc *application.TaskService, r *render.Renderer) *TasksHandler {
	return &TasksHandler{svc: svc, render: r}
}

func (h *TasksHandler) List(w http.ResponseWriter, r *http.Request) {
	f := parseTaskFilter(r)
	tasks, err := h.svc.List(f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	projects, _ := h.svc.ListActive()

	h.render.Page(w, http.StatusOK, "task_list.html", render.PageData{
		Title: "Tasks",
		Content: map[string]any{
			"Tasks":      tasks,
			"Projects":   projects,
			"Statuses":   domain.TaskStatuses,
			"Categories": domain.TaskCategories,
			"Filter":     f,
		},
	})
}

func (h *TasksHandler) TaskListPartial(w http.ResponseWriter, r *http.Request) {
	f := parseTaskFilter(r)
	tasks, err := h.svc.List(f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.render.Partial(w, http.StatusOK, "task-list-body", tasks)
}

func (h *TasksHandler) New(w http.ResponseWriter, r *http.Request) {
	projects, _ := h.svc.ListActive()
	preselect, _ := strconv.Atoi(r.URL.Query().Get("project_id"))
	meetingID, _ := strconv.Atoi(r.URL.Query().Get("meeting_id"))

	task := &domain.Task{
		ProjectID: preselect,
		Status:    "todo",
		Category:  "other",
	}
	if meetingID > 0 {
		task.MeetingID = &meetingID
	}

	h.render.Page(w, http.StatusOK, "task_form.html", render.PageData{
		Title: "New Task",
		Content: map[string]any{
			"Task":       task,
			"Projects":   projects,
			"Statuses":   domain.TaskStatuses,
			"Categories": domain.TaskCategories,
		},
	})
}

func (h *TasksHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	projectID, _ := strconv.Atoi(r.FormValue("project_id"))
	estimatedHours, _ := strconv.ParseFloat(r.FormValue("estimated_hours"), 64)
	isExternal := r.FormValue("is_external") == "1"

	t := &domain.Task{
		ProjectID:      projectID,
		Title:          r.FormValue("title"),
		StartDate:      r.FormValue("start_date"),
		Deadline:       r.FormValue("deadline"),
		EstimatedHours: estimatedHours,
		Status:         r.FormValue("status"),
		Category:       r.FormValue("category"),
		IsExternal:     isExternal,
		Description:    r.FormValue("description"),
	}

	if mid := r.FormValue("meeting_id"); mid != "" {
		id, _ := strconv.Atoi(mid)
		if id > 0 {
			t.MeetingID = &id
		}
	}

	projects, _ := h.svc.ListActive()
	errors := domain.ValidateTask(t)
	if len(errors) > 0 {
		h.render.Page(w, http.StatusUnprocessableEntity, "task_form.html", render.PageData{
			Title: "New Task",
			Content: map[string]any{
				"Task":       t,
				"Projects":   projects,
				"Statuses":   domain.TaskStatuses,
				"Categories": domain.TaskCategories,
				"Errors":     errors,
			},
			Flash: "Please fix the errors below.",
		})
		return
	}

	if err := h.svc.Create(t); err != nil {
		h.render.Page(w, http.StatusUnprocessableEntity, "task_form.html", render.PageData{
			Title: "New Task",
			Content: map[string]any{
				"Task":       t,
				"Projects":   projects,
				"Statuses":   domain.TaskStatuses,
				"Categories": domain.TaskCategories,
			},
			Flash: "Error: " + err.Error(),
		})
		return
	}

	if t.MeetingID != nil {
		http.Redirect(w, r, "/meetings/"+strconv.Itoa(*t.MeetingID), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/tasks", http.StatusSeeOther)
}

func (h *TasksHandler) Edit(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	task, err := h.svc.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	projects, _ := h.svc.ListActive()

	h.render.Page(w, http.StatusOK, "task_form.html", render.PageData{
		Title: "Edit Task",
		Content: map[string]any{
			"Task":       task,
			"Projects":   projects,
			"Statuses":   domain.TaskStatuses,
			"Categories": domain.TaskCategories,
		},
	})
}

func (h *TasksHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	projectID, _ := strconv.Atoi(r.FormValue("project_id"))
	estimatedHours, _ := strconv.ParseFloat(r.FormValue("estimated_hours"), 64)
	isExternal := r.FormValue("is_external") == "1"

	t := &domain.Task{
		ID:             id,
		ProjectID:      projectID,
		Title:          r.FormValue("title"),
		StartDate:      r.FormValue("start_date"),
		Deadline:       r.FormValue("deadline"),
		EstimatedHours: estimatedHours,
		Status:         r.FormValue("status"),
		Category:       r.FormValue("category"),
		IsExternal:     isExternal,
		Description:    r.FormValue("description"),
	}

	if mid := r.FormValue("meeting_id"); mid != "" {
		midInt, _ := strconv.Atoi(mid)
		if midInt > 0 {
			t.MeetingID = &midInt
		}
	}

	projects, _ := h.svc.ListActive()
	errors := domain.ValidateTask(t)
	if len(errors) > 0 {
		h.render.Page(w, http.StatusUnprocessableEntity, "task_form.html", render.PageData{
			Title: "Edit Task",
			Content: map[string]any{
				"Task":       t,
				"Projects":   projects,
				"Statuses":   domain.TaskStatuses,
				"Categories": domain.TaskCategories,
				"Errors":     errors,
			},
			Flash: "Please fix the errors below.",
		})
		return
	}

	if err := h.svc.Update(t); err != nil {
		h.render.Page(w, http.StatusUnprocessableEntity, "task_form.html", render.PageData{
			Title: "Edit Task",
			Content: map[string]any{
				"Task":       t,
				"Projects":   projects,
				"Statuses":   domain.TaskStatuses,
				"Categories": domain.TaskCategories,
			},
			Flash: "Error: " + err.Error(),
		})
		return
	}

	http.Redirect(w, r, "/tasks", http.StatusSeeOther)
}

func (h *TasksHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
		w.Header().Set("HX-Redirect", "/tasks")
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/tasks", http.StatusSeeOther)
}

func (h *TasksHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	status := r.FormValue("status")
	if err := h.svc.UpdateStatus(id, status); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	task, err := h.svc.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.render.Partial(w, http.StatusOK, "task-row", task)
}

func parseTaskFilter(r *http.Request) domain.TaskFilter {
	projectID, _ := strconv.Atoi(r.URL.Query().Get("project_id"))
	return domain.TaskFilter{
		ProjectID: projectID,
		Status:    r.URL.Query().Get("status"),
		Category:  r.URL.Query().Get("category"),
		Overdue:   r.URL.Query().Get("overdue") == "1",
	}
}
