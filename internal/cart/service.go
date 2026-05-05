package cart

type Service struct {
    repo *Repository
}

func NewService(repo *Repository) *Service {
    return &Service{repo}
}

func (s *Service) AddToCart(userID, productID uint, quantity int) error {
    cart, err := s.repo.GetCartByUser(userID)

    // if cart doesn't exist → create it
    if err != nil {
        cart = &Cart{UserID: userID}
        if err := s.repo.CreateCart(cart); err != nil {
            return err
        }
    }

    // check if product already exists
    for _, item := range cart.Items {
        if item.ProductID == productID {
            item.Quantity += quantity
            return s.repo.UpdateItem(&item)
        }
    }

    // create new item
    newItem := &CartItem{
        CartID:    cart.ID,
        ProductID: productID,
        Quantity:  quantity,
    }

    return s.repo.AddItem(newItem)
}