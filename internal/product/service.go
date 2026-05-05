package product

type Service struct {
    repo *Repository
}

func NewService(repo *Repository) *Service {
    return &Service{repo}
}

func (s *Service) CreateProduct(product *Product) error {
    return s.repo.Create(product)
}

func (s *Service) GetProducts() ([]Product, error) {
    return s.repo.FindAll()
}