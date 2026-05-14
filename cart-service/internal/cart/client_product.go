package cart

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type ProductClient struct {
	BaseURL string
}

func NewProductClient(baseURL string) *ProductClient {
	return &ProductClient{BaseURL: baseURL}
}

func (p *ProductClient) GetProducts(
	productIDs []uint,
	authorizationHeader string,
) (map[uint]ProductDTO, error) {

	url := fmt.Sprintf("%s/api/v1/buyer/products/bulk", p.BaseURL)

	payload := map[string]interface{}{
		"product_ids": productIDs,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authorizationHeader)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to fetch products")
	}

	var result BulkProductsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	productsMap := make(map[uint]ProductDTO)
	for _, product := range result.Data {
		productsMap[product.ID] = product
	}

	return productsMap, nil
}