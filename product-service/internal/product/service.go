package product

import (
	"errors"
	"math"
	"product-service/pkg/utils"

	"gorm.io/datatypes"
)

type Service struct {
	Repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{Repo: repo}
}

func (s *Service) CreateProduct(
	sellerID uint,
	req CreateProductRequest,
) (*ProductResponse, error) {

	product := Product{
		SellerID:        sellerID,
		Title:           req.Title,
		Description:     req.Description,
		Price:           req.Price,
		Stock:           req.Stock,
		Category:        req.Category,
		ImageURLs:       datatypes.JSONSlice[string](req.ImageURLs),
		Offers:          req.Offers,
		ReturnAvailable: req.ReturnAvailable,
		Warranty:        req.Warranty,
		IsActive:        true,
	}

	if err := s.Repo.Create(&product); err != nil {
		return nil, errors.New("failed to create product")
	}

	resp := productToResponse(product)
	return &resp, nil
}

func (s *Service) GetAllProducts(page int) (*BuyerProductsPaginatedResponse, error) {
	const limit = 5

	total, err := s.Repo.CountActive()
	if err != nil {
		return nil, errors.New("failed to count products")
	}

	totalPages := getTotalPages(total, limit)
	page = normalizePage(page, totalPages)
	offset := (page - 1) * limit

	products, err := s.Repo.FindActiveOrderedPaginated(limit, offset)
	if err != nil {
		return nil, errors.New("failed to fetch products")
	}

	response := make([]BuyerProductResponse, 0, len(products))
	for _, p := range products {
		image := ""
		if len(p.ImageURLs) > 0 {
			image = p.ImageURLs[0]
		}

		response = append(response, BuyerProductResponse{
			ID:               p.ID,
			Title:            p.Title,
			ImageURL:         image,
			Description:      p.Description,
			Price:            p.Price,
			Offers:           p.Offers,
			ExpectedDelivery: utils.GetExpectedDelivery(),
		})
	}

	return &BuyerProductsPaginatedResponse{
		Products: response,
		Pagination: PaginationMeta{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *Service) GetProductByID(productID uint) (*BuyerProductDetailResponse, error) {
	p, err := s.Repo.FirstActiveByID(productID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	return &BuyerProductDetailResponse{
		ID:               p.ID,
		Title:            p.Title,
		Description:      p.Description,
		Price:            p.Price,
		Category:         p.Category,
		ImageURLs:        []string(p.ImageURLs),
		Offers:           p.Offers,
		Rating:           p.Rating,
		ReturnAvailable:  p.ReturnAvailable,
		Warranty:         p.Warranty,
		InStock:          p.Stock > 0,
		Stock:            p.Stock,
		ExpectedDelivery: utils.GetExpectedDelivery(),
	}, nil
}

func (s *Service) GetSellerProductByID(
	sellerID uint,
	productID uint,
) (*ProductResponse, error) {

	p, err := s.Repo.FirstBySellerAndProductID(sellerID, productID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	resp := productToResponse(*p)
	return &resp, nil
}

func (s *Service) GetSellerProducts(
	sellerID uint,
	page int,
) (*SellerProductsPaginatedResponse, error) {

	const limit = 5

	total, err := s.Repo.CountBySeller(sellerID)
	if err != nil {
		return nil, errors.New("failed to count seller products")
	}

	totalPages := getTotalPages(total, limit)
	page = normalizePage(page, totalPages)
	offset := (page - 1) * limit

	products, err := s.Repo.FindBySellerPaginated(sellerID, limit, offset)
	if err != nil {
		return nil, errors.New("failed to fetch seller products")
	}

	out := make([]ProductResponse, 0, len(products))
	for _, p := range products {
		out = append(out, productToResponse(p))
	}

	return &SellerProductsPaginatedResponse{
		Products: out,
		Pagination: PaginationMeta{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

func getTotalPages(total int64, limit int) int {
	if total == 0 {
		return 0
	}
	return int(math.Ceil(float64(total) / float64(limit)))
}

func normalizePage(page int, totalPages int) int {
	if page < 1 {
		return 1
	}
	if totalPages > 0 && page > totalPages {
		return totalPages
	}
	return page
}

func (s *Service) UpdateProduct(
	sellerID uint,
	productID uint,
	req UpdateProductRequest,
) (*ProductResponse, error) {

	product, err := s.Repo.FirstBySellerAndProductID(sellerID, productID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	// Ownership is established; normalize/validate each provided field in the
	// same order as before the repository refactor (load first, then validate).

	if req.Title != nil {
		title := utils.NormalizeTitle(*req.Title)
		if err := utils.ValidateTitle(title); err != nil {
			return nil, err
		}
		product.Title = title
	}

	if req.Description != nil {
		description := utils.NormalizeDescription(*req.Description)
		if err := utils.ValidateDescription(description); err != nil {
			return nil, err
		}
		product.Description = description
	}

	if req.Price != nil {
		if err := utils.ValidatePrice(*req.Price); err != nil {
			return nil, err
		}
		product.Price = *req.Price
	}

	if req.Stock != nil {
		if err := utils.ValidateStock(*req.Stock); err != nil {
			return nil, err
		}
		product.Stock = *req.Stock
	}

	if req.Category != nil {
		category := utils.NormalizeCategory(*req.Category)
		if err := utils.ValidateCategory(category); err != nil {
			return nil, err
		}
		product.Category = category
	}

	if req.ImageURLs != nil {
		imageURLs := utils.NormalizeStringSlice(*req.ImageURLs)
		if err := utils.ValidateImageURLs(imageURLs); err != nil {
			return nil, err
		}
		product.ImageURLs = datatypes.JSONSlice[string](imageURLs)
	}

	if req.Offers != nil {
		offers := utils.NormalizeOptionalText(*req.Offers)
		if err := utils.ValidateOffers(offers); err != nil {
			return nil, err
		}
		product.Offers = offers
	}

	if req.ReturnAvailable != nil {
		product.ReturnAvailable = *req.ReturnAvailable
	}

	if req.Warranty != nil {
		warranty := utils.NormalizeOptionalText(*req.Warranty)
		if err := utils.ValidateWarranty(warranty); err != nil {
			return nil, err
		}
		product.Warranty = warranty
	}

	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}

	if err := s.Repo.Save(product); err != nil {
		return nil, errors.New("failed to update product")
	}

	resp := productToResponse(*product)
	return &resp, nil
}

func (s *Service) DeleteProduct(sellerID uint, productID uint) error {
	product, err := s.Repo.FirstBySellerAndProductID(sellerID, productID)
	if err != nil {
		return errors.New("product not found")
	}

	product.IsActive = false
	if err := s.Repo.Save(product); err != nil {
		return errors.New("failed to deactivate product")
	}

	return nil
}

func (s *Service) GetProductsByIDs(ids []uint) ([]ProductResponse, error) {
	products, err := s.Repo.FindActiveByIDs(ids)
	if err != nil {
		return nil, errors.New("failed to fetch products")
	}

	out := make([]ProductResponse, 0, len(products))
	for _, p := range products {
		out = append(out, productToResponse(p))
	}
	return out, nil
}
func (s *Service) GetInventory(
	productID uint,
) (*InventoryResponse, error) {

	product, err := s.Repo.FindProductByID(productID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	return &InventoryResponse{
		ProductID: product.ID,
		Stock:     product.Stock,
	}, nil
}
func (s *Service) ReduceStock(
	productID uint,
	quantity int,
) error {

	if quantity <= 0 {
		return errors.New("invalid quantity")
	}

	err := s.Repo.ReduceStock(
		productID,
		quantity,
	)

	if err != nil {
		return errors.New("insufficient stock")
	}

	return nil
}