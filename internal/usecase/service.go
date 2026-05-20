package usecase

import (
	"context"
	"my-go-server/internal/domain"
)

type Service struct{
	repository domain.Repository
}

func NewService(repository domain.Repository) domain.UseCase {
	return &Service{
		repository: repository,
	}
}
	
func (s *Service) GetMessage(ctx context.Context) (string, error) {
		return s.repository.GetMessage(ctx)
	}

func (s *Service) SaveMessage(ctx context.Context, value string) (int, error) {
		return s.repository.SaveMessage(ctx, value)
	}