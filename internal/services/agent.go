package services

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

// NoOpAgent is the default stub implementation
type NoOpAgent struct{}

func NewNoOpAgent() *NoOpAgent {
	return &NoOpAgent{}
}

func (a *NoOpAgent) ExtractTasksFromMeeting(ctx context.Context, notes string, project string) ([]ExtractedTask, error) {
	return nil, nil
}

func (a *NoOpAgent) SummarizeProjectStatus(ctx context.Context, projectID int) (string, error) {
	return "", nil
}
