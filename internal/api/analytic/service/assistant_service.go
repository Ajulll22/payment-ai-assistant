package service

import "context"

type AssistantService interface {
}

type assistantService struct {
}

func (s *assistantService) ChatAssistant(ctx context.Context, roomID, prompt string) error {

	return nil
}
