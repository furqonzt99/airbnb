package rating

type RatingResponse struct {
	HouseID  int    `json:"house_id"`
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Rating   int    `json:"rating"`
	Comment  string `json:"comment"`
}
