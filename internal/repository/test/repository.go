package repository

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) GetMessage() string {
		return "znayu ya chto porvetsa pizdyukha!"
	}