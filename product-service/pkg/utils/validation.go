package utils

import (
	"errors"
	"net/url"
	"regexp"
	"strings"
	"unicode"
)

var AllowedCategories = map[string]bool{
	"Electronics":    true,
	"Fashion":        true,
	"Home & Kitchen": true,
	"Books":          true,
	"Beauty":         true,
	"Toys":           true,
	"Sports":         true,
	"Groceries":      true,
	"Furniture":      true,
	"Automotive":     true,
}

// -------------------------
// NORMALIZERS
// -------------------------

func NormalizeTitle(title string) string {
	return strings.TrimSpace(title)
}

func NormalizeCategory(category string) string {
	return strings.TrimSpace(category)
}

func NormalizeDescription(desc string) string {
	return strings.TrimSpace(desc)
}

func NormalizeOptionalText(value string) string {
	return strings.TrimSpace(value)
}

func NormalizeStringSlice(values []string) []string {
	normalized := make([]string, len(values))

	for i, value := range values {
		normalized[i] = strings.TrimSpace(value)
	}

	return normalized
}

func ValidateID(id int, field string) error {
	if id <= 0 {
		return errors.New(field + " must be greater than 0")
	}

	return nil
}

// -------------------------
// TITLE
// -------------------------

func ValidateTitle(title string) error {

	title = strings.TrimSpace(title)

	if title == "" {
		return errors.New("title is required")
	}

	if len(title) < 3 {
		return errors.New("title must be at least 3 characters")
	}

	if len(title) > 120 {
		return errors.New("title cannot exceed 120 characters")
	}

	hasLetter := false

	for _, c := range title {
		if unicode.IsLetter(c) {
			hasLetter = true
			break
		}
	}

	if !hasLetter {
		return errors.New("title must contain meaningful text")
	}

	regex := `^[a-zA-Z0-9\s\-\&\(\)\+\.]+$`

	matched, _ := regexp.MatchString(regex, title)

	if !matched {
		return errors.New("title contains invalid characters")
	}

	return nil
}

// -------------------------
// DESCRIPTION
// -------------------------

func ValidateDescription(description string) error {

	description = strings.TrimSpace(description)

	if description == "" {
		return errors.New("description is required")
	}

	if len(description) < 20 {
		return errors.New("description too short")
	}

	if len(description) > 2500 {
		return errors.New("description cannot exceed 2500 characters")
	}

	return nil
}

// -------------------------
// PRICE
// -------------------------

func ValidatePrice(price float64) error {

	if price <= 0 {
		return errors.New("price must be greater than 0")
	}

	if price > 10000000 {
		return errors.New("price too large")
	}

	return nil
}

// -------------------------
// STOCK
// -------------------------

func ValidateStock(stock int) error {

	if stock < 0 {
		return errors.New("stock cannot be negative")
	}

	if stock > 100000 {
		return errors.New("stock value too large")
	}

	return nil
}

// -------------------------
// CATEGORY
// -------------------------

func ValidateCategory(category string) error {

	category = strings.TrimSpace(category)

	if category == "" {
		return errors.New("category is required")
	}

	if !AllowedCategories[category] {
		return errors.New("invalid category")
	}

	return nil
}

// -------------------------
// IMAGE URLS
// -------------------------

func ValidateImageURLs(images []string) error {

	if len(images) == 0 {
		return errors.New("at least one image is required")
	}

	if len(images) > 8 {
		return errors.New("maximum 8 images allowed")
	}

	for _, image := range images {

		image = strings.TrimSpace(image)

		if image == "" {
			return errors.New("image url cannot be empty")
		}

		parsed, err := url.ParseRequestURI(image)

		if err != nil || parsed.Scheme == "" || parsed.Host == "" {
			return errors.New("invalid image url format")
		}
	}

	return nil
}

// -------------------------
// OFFERS
// -------------------------

func ValidateOffers(offers string) error {

	offers = strings.TrimSpace(offers)

	if offers == "" {
		return nil
	}

	if len(offers) > 120 {
		return errors.New("offers text too long")
	}

	return nil
}

// -------------------------
// WARRANTY
// -------------------------

func ValidateWarranty(warranty string) error {

	warranty = strings.TrimSpace(warranty)

	if warranty == "" {
		return nil
	}

	if len(warranty) > 120 {
		return errors.New("warranty text too long")
	}

	return nil
}
