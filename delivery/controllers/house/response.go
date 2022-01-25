package house

import "github.com/furqonzt99/airbnb/model"

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
	ID       uint    `json:"id"`
	UserID   uint    `json:"userid"`
	UserName string  `json:"username"`
	Title    string  `json:"title"`
	Address  string  `json:"address"`
	City     string  `json:"city"`
	Price    float64 `json:"price"`
	Features []model.Feature
}
