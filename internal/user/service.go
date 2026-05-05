package user

type Service struct {
    repo *Repository
}

func NewService(repo *Repository) *Service {
    return &Service{repo}
}

func (s *Service) CreateUser(user *User) error {
    return s.repo.Create(user)
}