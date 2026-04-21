package test

type Service struct{
	repository MessageRepository
}

func NewService(repository MessageRepository) *Service {
	return &Service{
		repository: repository,
	}
}
	
func (s *Service) GetMessage() string {
		return s.repository.GetMessage()
	}