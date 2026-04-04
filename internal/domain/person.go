package domain

type Person struct {
	ID        int
	Name      string
	Company   string
	Role      string
	Email     string
	Phone     string
	CreatedAt string
	UpdatedAt string
}

type PersonRepository interface {
	List() ([]Person, error)
	GetByID(id int) (*Person, error)
	Create(p *Person) error
	Update(p *Person) error
	Delete(id int) error
}
