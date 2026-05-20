package usecase

import (
	"context"
	"fmt"
	"my-go-server/internal/domain"
	"my-go-server/internal/models"

	"golang.org/x/crypto/bcrypt"
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

func (s *Service) RegisterUser(ctx context.Context, login, password string) (*models.User, error) {
	if login == "" {
		return nil, fmt.Errorf("login is required")
	}

	if password == "" {
		return nil, fmt.Errorf("password is required")
	}
	
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	existing, err := s.repository.GetUserByLogin(ctx, login)

	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}

	if existing != nil {
		return nil, fmt.Errorf("user with login %s already exists", login)
	}
	return s.repository.RegisterUser(ctx, login, string(hashedPassword))
}

func (s *Service) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	return s.repository.GetUserByLogin(ctx, login)
}

func (s *Service) CreateSession(ctx context.Context, userID string) (*models.Session, error) {
	return s.repository.CreateSession(ctx, userID)
}