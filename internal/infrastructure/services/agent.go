package services

import (
	"context"

	"github.com/daniel/ppm/internal/domain"
)

// NoOpAgent is the default stub implementation of domain.AgentService
type NoOpAgent struct{}

func NewNoOpAgent() *NoOpAgent {
	return &NoOpAgent{}
}

func (a *NoOpAgent) ExtractTasksFromMeeting(ctx context.Context, notes string, project string) ([]domain.ExtractedTask, error) {
	return nil, nil
}

func (a *NoOpAgent) SummarizeProjectStatus(ctx context.Context, projectID int) (string, error) {
	return "", nil
}
