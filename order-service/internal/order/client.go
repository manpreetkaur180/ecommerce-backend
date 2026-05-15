package order

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
	"bytes"
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

type CartItemResponse struct {
	ProductID uint `json:"product_id"`
	SellerID  uint `json:"seller_id"`

	ProductName  string `json:"product_name"`
	ProductImage string `json:"product_image"`

	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

type AddressResponse struct {
	ID uint `json:"id"`

	FullName string `json:"full_name"`
	Phone    string `json:"phone"`

	AddressLine1 string `json:"address_line1"`
	AddressLine2 string `json:"address_line2"`
	Landmark     string `json:"landmark"`
	City         string `json:"city"`
	State        string `json:"state"`
	Country      string `json:"country"`
	Pincode      string `json:"pincode"`
}

type InventoryResponse struct {
	ProductID uint `json:"product_id"`
	Stock     int  `json:"stock"`
}

func FetchCartItems(token string) ([]CartItemResponse, error) {

	url := fmt.Sprintf(
		"%s/api/v1/cart",
		os.Getenv("CART_SERVICE_URL"),
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var result struct {
		Data []CartItemResponse `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

func FetchUserAddresses(
	token string,
) ([]AddressResponse, error) {

	url := fmt.Sprintf(
		"%s/api/v1/user/addresses",
		os.Getenv("USER_SERVICE_URL"),
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var result struct {
		Data []AddressResponse `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}
func FetchInventory(
	productID uint,
) (*InventoryResponse, error) {

	url := fmt.Sprintf(
		"%s/internal/products/%d/inventory",
		os.Getenv("PRODUCT_SERVICE_URL"),
		productID,
	)

	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var result struct {
		Data InventoryResponse `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result.Data, nil
}

func ReduceInventory(
	productID uint,
	quantity int,
) error {

	url := fmt.Sprintf(
		"%s/internal/products/%d/reduce-stock",
		os.Getenv("PRODUCT_SERVICE_URL"),
		productID,
	)

	body := map[string]int{
		"quantity": quantity,
	}

	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequest(
		"PATCH",
		url,
		bytes.NewBuffer(jsonBody),
	)

	if err != nil {
		return err
	}

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("failed to reduce stock")
	}

	return nil
}