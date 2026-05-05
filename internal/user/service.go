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
func (s *Service) LoginByEmail(email string) (*User, error) {
    user, err := s.repo.FindByEmail(email)
    if err != nil {
        return nil, err
    }

    return user, nil
}