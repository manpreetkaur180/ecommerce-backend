package cart

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ProductDTO struct {
	ID          uint     `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	Stock       int      `json:"stock"`
	ImageURLs   []string `json:"image_urls"`
}

type ProductClient struct {
	BaseURL string
}

func NewProductClient(baseURL string) *ProductClient {
	return &ProductClient{BaseURL: baseURL}
}

func (p *ProductClient) GetProduct(productID uint, authorizationHeader string) (*ProductDTO, error) {

	url := fmt.Sprintf("%s/api/v1/buyer/products/%d", p.BaseURL, productID)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", authorizationHeader)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("product not found")
	}

	var result struct {
		Data ProductDTO `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result.Data, nil
}
