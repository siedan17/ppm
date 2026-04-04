package http

import (
	"net/http"
	"strconv"

	"github.com/daniel/ppm/internal/application"
	"github.com/daniel/ppm/internal/delivery/render"
	"github.com/daniel/ppm/internal/domain"
)

type MeetingsHandler struct {
	svc    *application.MeetingService
	render *render.Renderer
}

func NewMeetingsHandler(svc *application.MeetingService, r *render.Renderer) *MeetingsHandler {
	return &MeetingsHandler{svc: svc, render: r}
}

func (h *MeetingsHandler) List(w http.ResponseWriter, r *http.Request) {
	f := domain.MeetingFilter{
		DateFrom: r.URL.Query().Get("date_from"),
		DateTo:   r.URL.Query().Get("date_to"),
	}
	if pid := r.URL.Query().Get("project_id"); pid != "" {
		f.ProjectID, _ = strconv.Atoi(pid)
	}

	meetings, err := h.svc.List(f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	projects, _ := h.svc.ListActive()

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
	projects, _ := h.svc.ListActive()
	preselect, _ := strconv.Atoi(r.URL.Query().Get("project_id"))

	h.render.Page(w, http.StatusOK, "meeting_form.html", render.PageData{
		Title: "New Meeting",
		Content: map[string]any{
			"Meeting":      &domain.Meeting{ProjectID: preselect, MeetingType: "external"},
			"Projects":     projects,
			"MeetingTypes": domain.MeetingTypes,
		},
	})
}

func (h *MeetingsHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	projectID, _ := strconv.Atoi(r.FormValue("project_id"))
	m := &domain.Meeting{
		ProjectID:   projectID,
		Date:        r.FormValue("date"),
		MeetingType: r.FormValue("meeting_type"),
		Title:       r.FormValue("title"),
		Notes:       r.FormValue("notes"),
	}

	projects, _ := h.svc.ListActive()
	errors := domain.ValidateMeeting(m)
	if len(errors) > 0 {
		h.render.Page(w, http.StatusUnprocessableEntity, "meeting_form.html", render.PageData{
			Title: "New Meeting",
			Content: map[string]any{
				"Meeting":      m,
				"Projects":     projects,
				"MeetingTypes": domain.MeetingTypes,
				"Errors":       errors,
			},
			Flash: "Please fix the errors below.",
		})
		return
	}

	if err := h.svc.Create(m); err != nil {
		h.render.Page(w, http.StatusUnprocessableEntity, "meeting_form.html", render.PageData{
			Title: "New Meeting",
			Content: map[string]any{
				"Meeting":      m,
				"Projects":     projects,
				"MeetingTypes": domain.MeetingTypes,
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

	detail, err := h.svc.GetDetail(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	h.render.Page(w, http.StatusOK, "meeting_detail.html", render.PageData{
		Title: detail.Meeting.Title,
		Content: map[string]any{
			"Meeting":   detail.Meeting,
			"Tasks":     detail.Tasks,
			"AllPeople": detail.AllPeople,
		},
	})
}

func (h *MeetingsHandler) Edit(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	meeting, err := h.svc.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	projects, _ := h.svc.ListActive()

	h.render.Page(w, http.StatusOK, "meeting_form.html", render.PageData{
		Title: "Edit Meeting",
		Content: map[string]any{
			"Meeting":      meeting,
			"Projects":     projects,
			"MeetingTypes": domain.MeetingTypes,
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
	m := &domain.Meeting{
		ID:          id,
		ProjectID:   projectID,
		Date:        r.FormValue("date"),
		MeetingType: r.FormValue("meeting_type"),
		Title:       r.FormValue("title"),
		Notes:       r.FormValue("notes"),
	}

	projects, _ := h.svc.ListActive()
	errors := domain.ValidateMeeting(m)
	if len(errors) > 0 {
		h.render.Page(w, http.StatusUnprocessableEntity, "meeting_form.html", render.PageData{
			Title: "Edit Meeting",
			Content: map[string]any{
				"Meeting":      m,
				"Projects":     projects,
				"MeetingTypes": domain.MeetingTypes,
				"Errors":       errors,
			},
			Flash: "Please fix the errors below.",
		})
		return
	}

	if err := h.svc.Update(m); err != nil {
		h.render.Page(w, http.StatusUnprocessableEntity, "meeting_form.html", render.PageData{
			Title: "Edit Meeting",
			Content: map[string]any{
				"Meeting":      m,
				"Projects":     projects,
				"MeetingTypes": domain.MeetingTypes,
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

	if err := h.svc.Delete(id); err != nil {
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
	if err := h.svc.AddParticipant(meetingID, personID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/meetings/"+strconv.Itoa(meetingID), http.StatusSeeOther)
}

func (h *MeetingsHandler) RemoveParticipant(w http.ResponseWriter, r *http.Request) {
	meetingID, _ := strconv.Atoi(r.PathValue("id"))
	personID, _ := strconv.Atoi(r.PathValue("pid"))

	if err := h.svc.RemoveParticipant(meetingID, personID); err != nil {
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

	meeting, err := h.svc.GetByID(meetingID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	projects, _ := h.svc.ListActive()

	task := &domain.Task{
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
			"Statuses":   domain.TaskStatuses,
			"Categories": domain.TaskCategories,
			"Meeting":    meeting,
		},
	})
}
