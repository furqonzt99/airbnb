package common

type CallbackRequest struct {
	ExternalID string `json:"external_id"`
	PaymentMethod string `json:"payment_method"`
	BankID string `json:"bank_code"`
	PaidAt string `json:"paid_at"`
	Status string `json:"status"`
}