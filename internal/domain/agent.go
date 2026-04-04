package domain

import "context"

type ExtractedTask struct {
	Title          string
	Description    string
	Category       string
	EstimatedHours float64
}

type AgentService interface {
	ExtractTasksFromMeeting(ctx context.Context, notes string, project string) ([]ExtractedTask, error)
	SummarizeProjectStatus(ctx context.Context, projectID int) (string, error)
}
