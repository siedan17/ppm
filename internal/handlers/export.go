package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/daniel/ppm/internal/repository"
)

type ExportHandler struct {
	meetingRepo *repository.MeetingRepo
	projectRepo *repository.ProjectRepo
}

func NewExportHandler(meetingRepo *repository.MeetingRepo, projectRepo *repository.ProjectRepo) *ExportHandler {
	return &ExportHandler{meetingRepo: meetingRepo, projectRepo: projectRepo}
}

func (h *ExportHandler) ExportMeeting(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	meeting, err := h.meetingRepo.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# %s\n\n", meeting.Title))
	sb.WriteString(fmt.Sprintf("**Date:** %s\n", meeting.Date))
	sb.WriteString(fmt.Sprintf("**Project:** %s\n", meeting.ProjectName))
	sb.WriteString(fmt.Sprintf("**Type:** %s\n\n", meeting.MeetingType))

	if len(meeting.Participants) > 0 {
		sb.WriteString("**Participants:**\n")
		for _, p := range meeting.Participants {
			sb.WriteString(fmt.Sprintf("- %s", p.Name))
			if p.Company != "" {
				sb.WriteString(fmt.Sprintf(" (%s)", p.Company))
			}
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	}

	sb.WriteString("## Notes\n\n")
	sb.WriteString(meeting.Notes)
	sb.WriteString("\n")

	filename := fmt.Sprintf("%s_%s.md", meeting.Date, sanitizeFilename(meeting.Title))
	w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	w.Write([]byte(sb.String()))
}

func (h *ExportHandler) ExportProjectMeetings(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	project, err := h.projectRepo.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	meetings, err := h.meetingRepo.ListByProject(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# %s — Meeting Notes\n\n", project.Name))

	for _, m := range meetings {
		sb.WriteString(fmt.Sprintf("## %s — %s\n\n", m.Date, m.Title))
		sb.WriteString(fmt.Sprintf("**Type:** %s\n\n", m.MeetingType))
		sb.WriteString(m.Notes)
		sb.WriteString("\n\n---\n\n")
	}

	filename := fmt.Sprintf("%s_meetings.md", sanitizeFilename(project.Name))
	w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	w.Write([]byte(sb.String()))
}

func sanitizeFilename(s string) string {
	replacer := strings.NewReplacer(" ", "_", "/", "_", "\\", "_", ":", "_")
	return strings.ToLower(replacer.Replace(s))
}
