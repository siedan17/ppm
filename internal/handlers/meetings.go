package handlers

import (
	"net/http"
	"strconv"

	"github.com/daniel/ppm/internal/models"
	"github.com/daniel/ppm/internal/render"
	"github.com/daniel/ppm/internal/repository"
)

type MeetingsHandler struct {
	repo        *repository.MeetingRepo
	projectRepo *repository.ProjectRepo
	personRepo  *repository.PersonRepo
	taskRepo    *repository.TaskRepo
	render      *render.Renderer
}

func NewMeetingsHandler(repo *repository.MeetingRepo, projectRepo *repository.ProjectRepo, personRepo *repository.PersonRepo, taskRepo *repository.TaskRepo, r *render.Renderer) *MeetingsHandler {
	return &MeetingsHandler{repo: repo, projectRepo: projectRepo, personRepo: personRepo, taskRepo: taskRepo, render: r}
}

func (h *MeetingsHandler) List(w http.ResponseWriter, r *http.Request) {
	f := repository.MeetingFilter{
		DateFrom: r.URL.Query().Get("date_from"),
		DateTo:   r.URL.Query().Get("date_to"),
	}
	if pid := r.URL.Query().Get("project_id"); pid != "" {
		f.ProjectID, _ = strconv.Atoi(pid)
	}

	meetings, err := h.repo.List(f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	projects, _ := h.projectRepo.ListActive()

	h.render.Page(w, http.StatusOK, "meeting_list.html", render.PageData{
		Title: "Meetings",
		Content: map[string]any{
			"Meetings": meetings,
			"Projects": projects,
			"Filter":   f,
		},
	})
}

func (h *MeetingsHandler) New(w http.ResponseWriter, r *http.Request) {
	projects, _ := h.projectRepo.ListActive()
	preselect, _ := strconv.Atoi(r.URL.Query().Get("project_id"))

	h.render.Page(w, http.StatusOK, "meeting_form.html", render.PageData{
		Title: "New Meeting",
		Content: map[string]any{
			"Meeting":      &models.Meeting{ProjectID: preselect, MeetingType: "external"},
			"Projects":     projects,
			"MeetingTypes": models.MeetingTypes,
		},
	})
}

func (h *MeetingsHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	projectID, _ := strconv.Atoi(r.FormValue("project_id"))
	m := &models.Meeting{
		ProjectID:   projectID,
		Date:        r.FormValue("date"),
		MeetingType: r.FormValue("meeting_type"),
		Title:       r.FormValue("title"),
		Notes:       r.FormValue("notes"),
	}

	projects, _ := h.projectRepo.ListActive()
	errors := validateMeeting(m)
	if len(errors) > 0 {
		h.render.Page(w, http.StatusUnprocessableEntity, "meeting_form.html", render.PageData{
			Title: "New Meeting",
			Content: map[string]any{
				"Meeting":      m,
				"Projects":     projects,
				"MeetingTypes": models.MeetingTypes,
				"Errors":       errors,
			},
			Flash: "Please fix the errors below.",
		})
		return
	}

	if err := h.repo.Create(m); err != nil {
		h.render.Page(w, http.StatusUnprocessableEntity, "meeting_form.html", render.PageData{
			Title: "New Meeting",
			Content: map[string]any{
				"Meeting":      m,
				"Projects":     projects,
				"MeetingTypes": models.MeetingTypes,
			},
			Flash: "Error: " + err.Error(),
		})
		return
	}

	http.Redirect(w, r, "/meetings/"+strconv.Itoa(m.ID), http.StatusSeeOther)
}

func (h *MeetingsHandler) Detail(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	meeting, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Render markdown notes
	meeting.NotesHTML = h.render.RenderMarkdown(meeting.Notes)

	tasks, _ := h.taskRepo.ListByMeeting(id)
	allPeople, _ := h.personRepo.List()

	h.render.Page(w, http.StatusOK, "meeting_detail.html", render.PageData{
		Title: meeting.Title,
		Content: map[string]any{
			"Meeting":   meeting,
			"Tasks":     tasks,
			"AllPeople": allPeople,
		},
	})
}

func (h *MeetingsHandler) Edit(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	meeting, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	projects, _ := h.projectRepo.ListActive()

	h.render.Page(w, http.StatusOK, "meeting_form.html", render.PageData{
		Title: "Edit Meeting",
		Content: map[string]any{
			"Meeting":      meeting,
			"Projects":     projects,
			"MeetingTypes": models.MeetingTypes,
		},
	})
}

func (h *MeetingsHandler) Update(w http.ResponseWriter, r *http.Request) {
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
	m := &models.Meeting{
		ID:          id,
		ProjectID:   projectID,
		Date:        r.FormValue("date"),
		MeetingType: r.FormValue("meeting_type"),
		Title:       r.FormValue("title"),
		Notes:       r.FormValue("notes"),
	}

	projects, _ := h.projectRepo.ListActive()
	errors := validateMeeting(m)
	if len(errors) > 0 {
		h.render.Page(w, http.StatusUnprocessableEntity, "meeting_form.html", render.PageData{
			Title: "Edit Meeting",
			Content: map[string]any{
				"Meeting":      m,
				"Projects":     projects,
				"MeetingTypes": models.MeetingTypes,
				"Errors":       errors,
			},
			Flash: "Please fix the errors below.",
		})
		return
	}

	if err := h.repo.Update(m); err != nil {
		h.render.Page(w, http.StatusUnprocessableEntity, "meeting_form.html", render.PageData{
			Title: "Edit Meeting",
			Content: map[string]any{
				"Meeting":      m,
				"Projects":     projects,
				"MeetingTypes": models.MeetingTypes,
			},
			Flash: "Error: " + err.Error(),
		})
		return
	}

	http.Redirect(w, r, "/meetings/"+strconv.Itoa(id), http.StatusSeeOther)
}

func (h *MeetingsHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
		w.Header().Set("HX-Redirect", "/meetings")
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/meetings", http.StatusSeeOther)
}

func (h *MeetingsHandler) AddParticipant(w http.ResponseWriter, r *http.Request) {
	meetingID, _ := strconv.Atoi(r.PathValue("id"))
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	personID, _ := strconv.Atoi(r.FormValue("person_id"))
	if err := h.repo.AddParticipant(meetingID, personID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/meetings/"+strconv.Itoa(meetingID), http.StatusSeeOther)
}

func (h *MeetingsHandler) RemoveParticipant(w http.ResponseWriter, r *http.Request) {
	meetingID, _ := strconv.Atoi(r.PathValue("id"))
	personID, _ := strconv.Atoi(r.PathValue("pid"))

	if err := h.repo.RemoveParticipant(meetingID, personID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/meetings/"+strconv.Itoa(meetingID), http.StatusSeeOther)
}

func (h *MeetingsHandler) CreateTaskFromMeeting(w http.ResponseWriter, r *http.Request) {
	meetingID, _ := strconv.Atoi(r.PathValue("id"))

	meeting, err := h.repo.GetByID(meetingID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	projects, _ := h.projectRepo.ListActive()

	task := &models.Task{
		ProjectID: meeting.ProjectID,
		MeetingID: &meetingID,
		Status:    "todo",
		Category:  "other",
	}

	h.render.Page(w, http.StatusOK, "task_form.html", render.PageData{
		Title: "New Task from Meeting: " + meeting.Title,
		Content: map[string]any{
			"Task":       task,
			"Projects":   projects,
			"Statuses":   models.TaskStatuses,
			"Categories": models.TaskCategories,
			"Meeting":    meeting,
		},
	})
}

func validateMeeting(m *models.Meeting) map[string]string {
	errors := make(map[string]string)
	if m.Title == "" {
		errors["title"] = "Title is required"
	}
	if m.Date == "" {
		errors["date"] = "Date is required"
	}
	if m.ProjectID == 0 {
		errors["project_id"] = "Project is required"
	}
	return errors
}
