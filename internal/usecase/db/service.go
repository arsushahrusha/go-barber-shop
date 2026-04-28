package db

type DBService struct {
	repository DBRepositoryMessage
}

func NewDBService(repo DBRepositoryMessage) *DBService {
	return &DBService{
		repository: repo,
	}
}

func (s *DBService) Save(value string) (int, error) {
	return s.repository.Save(value)
}