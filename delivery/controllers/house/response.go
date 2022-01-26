package house

import (
	"github.com/furqonzt99/airbnb/delivery/controllers/rating"
	"github.com/furqonzt99/airbnb/model"
)

type CreateHouseResponseFormat struct {
	Message string      `json:"message"`
	Data    model.House `json:"data"`
}

type GetAllHouseResponseFormat struct {
	Message string        `json:"message"`
	Data    []model.House `json:"data"`
}

type GetHouseResponseFormat struct {
	Message string      `json:"message"`
	Data    model.House `json:"data"`
}

type PutHouseResponseFormat struct {
	Message string      `json:"message"`
	Data    model.House `json:"data"`
}

type DeleteHouseResponseFormat struct {
	Message string `json:"message"`
}

type HouseResponse struct {
	ID       uint              `json:"id"`
	UserID   uint              `json:"user_id"`
	UserName string            `json:"user_name"`
	Title    string            `json:"title"`
	Address  string            `json:"address"`
	City     string            `json:"city"`
	Price    float64           `json:"price"`
	Rating   int			   `json:"rating"`
	Features []FeatureResponse `json:"features"`
	Ratings  []rating.RatingResponse `json:"ratings"`
}

type FeatureResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
