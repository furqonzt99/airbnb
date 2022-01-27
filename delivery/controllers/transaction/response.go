package transaction

type TransactionResponse struct {
	ID int `json:"id"`
	UserID int `json:"user_id"`
	HouseID int `json:"house_id"`
	HostID int `json:"host_id"`
	InvoiceID string `json:"invoice_id"`
	PaymentUrl string `json:"payment_url"`
	PaymentChannel string `json:"payment_channel"`
	PaymentMethod string `json:"payment_method"`
	PaidAt string `json:"paid_at"`
	CheckinDate string `json:"checkin_date"`
	CheckoutDate string `json:"checkout_date"`
	TotalPrice float64 `json:"total_price"`
	Status string `json:"status"`
}