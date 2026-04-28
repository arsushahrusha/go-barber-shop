package db

type DBRepositoryMessage interface {
	InitTable() error
	Save(value string) (int, error)
}

type DBServiceMessage interface {
	Save(value string) (int, error)
}