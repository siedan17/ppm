package application

import (
	"fmt"
	"strings"

	"github.com/daniel/ppm/internal/domain"
)

type ExportService struct {
	meetings domain.MeetingRepository
	projects domain.ProjectRepository
}

func NewExportService(meetings domain.MeetingRepository, projects domain.ProjectRepository) *ExportService {
	return &ExportService{meetings: meetings, projects: projects}
}

type ExportResult struct {
	Filename string
	Content  string
}

func (s *ExportService) ExportMeeting(id int) (*ExportResult, error) {
	meeting, err := s.meetings.GetByID(id)
	if err != nil {
		return nil, err
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
	return &ExportResult{Filename: filename, Content: sb.String()}, nil
}

func (s *ExportService) ExportProjectMeetings(id int) (*ExportResult, error) {
	project, err := s.projects.GetByID(id)
	if err != nil {
		return nil, err
	}

	meetings, err := s.meetings.ListByProject(id)
	if err != nil {
		return nil, err
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
	return &ExportResult{Filename: filename, Content: sb.String()}, nil
}

func sanitizeFilename(s string) string {
	replacer := strings.NewReplacer(" ", "_", "/", "_", "\\", "_", ":", "_")
	return strings.ToLower(replacer.Replace(s))
}
