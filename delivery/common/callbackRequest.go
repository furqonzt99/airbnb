package common

type CallbackRequest struct {
	ExternalID string `json:"external_id"`
	PaymentMethod string `json:"payment_method"`
	PaymentChannel string `json:"payment_channel"`
	PaidAt string `json:"paid_at"`
	Status string `json:"status"`
}