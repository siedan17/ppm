package application

import "github.com/daniel/ppm/internal/domain"

type PersonService struct {
	people domain.PersonRepository
}

func NewPersonService(people domain.PersonRepository) *PersonService {
	return &PersonService{people: people}
}

func (s *PersonService) List() ([]domain.Person, error) {
	return s.people.List()
}

func (s *PersonService) GetByID(id int) (*domain.Person, error) {
	return s.people.GetByID(id)
}

func (s *PersonService) Create(p *domain.Person) error {
	return s.people.Create(p)
}

func (s *PersonService) Update(p *domain.Person) error {
	return s.people.Update(p)
}

func (s *PersonService) Delete(id int) error {
	return s.people.Delete(id)
}
