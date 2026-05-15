package seller

import "time"

type ApplySellerRequest struct {
	BusinessName        string `json:"business_name"`
	BusinessDescription string `json:"business_description"`
	GSTIN               string `json:"gstin"`
	AadharNumber        string `json:"aadhar_number"`
}

type SellerApplicationListResponse struct {
	ID                  uint      `json:"id"`
	UserID              uint      `json:"user_id"`
	BusinessName        string    `json:"business_name"`
	BusinessDescription string    `json:"business_description"`
	GSTIN               string    `json:"gstin"`
	AadharNumber        string    `json:"aadhar_number"`
	Status              string    `json:"status"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type SellerApplicationDecisionRequest struct {
	AdminNote string `json:"admin_note"`
}

type SellerApplicationDecisionResponse struct {
	ID                  uint      `json:"id"`
	UserID              uint      `json:"user_id"`
	BusinessName        string    `json:"business_name"`
	BusinessDescription string    `json:"business_description"`
	GSTIN               string    `json:"gstin"`
	AadharNumber        string    `json:"aadhar_number"`
	Status              string    `json:"status"`
	AdminNote           string    `json:"admin_note"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

func ToSellerApplicationDecisionResponse(application *SellerApplication) SellerApplicationDecisionResponse {
	return SellerApplicationDecisionResponse{
		ID:                  application.ID,
		UserID:              application.UserID,
		BusinessName:        application.BusinessName,
		BusinessDescription: application.BusinessDescription,
		GSTIN:               application.GSTIN,
		AadharNumber:        application.AadharNumber,
		Status:              application.Status,
		AdminNote:           application.AdminNote,
		CreatedAt:           application.CreatedAt,
		UpdatedAt:           application.UpdatedAt,
	}
}
