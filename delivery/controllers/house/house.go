package house

import (
	"net/http"
	"strconv"

	"github.com/furqonzt99/airbnb/delivery/common"
	"github.com/furqonzt99/airbnb/delivery/controllers/rating"
	"github.com/furqonzt99/airbnb/delivery/middleware"
	"github.com/furqonzt99/airbnb/helper"
	"github.com/furqonzt99/airbnb/model"
	"github.com/furqonzt99/airbnb/repository/house"
	"github.com/labstack/echo/v4"
)

type HouseController struct {
	Repo house.HouseInterface
}

func NewHouseControllers(prorep house.HouseInterface) *HouseController {
	return &HouseController{Repo: prorep}
}

func (hc HouseController) CreateHouseController() echo.HandlerFunc {

	return func(c echo.Context) error {
		user, _ := middleware.ExtractTokenUser(c)

		newHouseReq := CreateHouseRequestFormat{}
		c.Bind(&newHouseReq)

		newHouse := model.House{
			UserID:  uint(user.UserID),
			Title:   newHouseReq.Title,
			Address: newHouseReq.Address,
			City:    newHouseReq.City,
			Price:   newHouseReq.Price,
		}

		house, err := hc.Repo.Create(newHouse)
		if err != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		for _, feature := range newHouseReq.Features {
			err := hc.Repo.HouseHasFeature(model.HouseHasFeatures{
				HouseID:   house.ID,
				FeatureID: uint(feature),
			})
			if err != nil {
				return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
			}
		}

		return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())
	}
}

func (hc HouseController) GetAllHouseController() echo.HandlerFunc {

	return func(c echo.Context) error {
		page, _ := strconv.Atoi(c.QueryParam("page"))
		perpage, _ := strconv.Atoi(c.QueryParam("perpage"))
		search := c.QueryParam("search")

		if page == 0 {
			page = 1
		}

		if perpage == 0 {
			perpage = 10
		}

		offset := (page - 1) * perpage

		houses, _ := hc.Repo.GetAll(offset, perpage, search)

		if len(houses) == 0 {
			return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
		}

		data := []HouseResponse{}
		for _, item := range houses {
			featuresData := []FeatureResponse{}
			for _, feature := range item.Features {
				featuresData = append(featuresData, FeatureResponse{
					ID:   feature.ID,
					Name: feature.Name,
				})
			}
			
			ratingData := []rating.RatingResponse{}
			ratings := []int{}

			for _, r := range item.Ratings {
				ratingData = append(ratingData, rating.RatingResponse{
					HouseID:  int(r.HouseID),
					UserID:   int(r.UserID),
					Username: r.User.Name,
					Rating:   r.Rating,
					Comment:  r.Comment,
				})

				ratings = append(ratings, r.Rating)
			}

			rating := helper.CalculateRatings(ratings)

			data = append(
				data, HouseResponse{
					ID:       item.ID,
					UserID:   item.User.ID,
					UserName: item.User.Name,
					Title:    item.Title,
					Address:  item.Address,
					City:     item.City,
					Price:    item.Price,
					Rating: rating,
					Features: featuresData,
					Ratings: ratingData,
				},
			)
		}
		return c.JSON(http.StatusOK, common.PaginationResponse(page, perpage, data))
	}
}

func (hc HouseController) GetMyHouseController() echo.HandlerFunc {

	return func(c echo.Context) error {

		user, _ := middleware.ExtractTokenUser(c)

		house, _ := hc.Repo.GetAllMine(user.UserID)

		data := []HouseResponse{}
		for _, item := range house {
			featuresData := []FeatureResponse{}
			for _, feature := range item.Features {
				featuresData = append(featuresData, FeatureResponse{
					ID:   feature.ID,
					Name: feature.Name,
				})
			}

			ratingData := []rating.RatingResponse{}
			ratings := []int{}

			for _, r := range item.Ratings {
				ratingData = append(ratingData, rating.RatingResponse{
					HouseID:  int(r.HouseID),
					UserID:   int(r.UserID),
					Username: r.User.Name,
					Rating:   r.Rating,
					Comment:  r.Comment,
				})

				ratings = append(ratings, r.Rating)
			}

			rating := helper.CalculateRatings(ratings)

			data = append(
				data, HouseResponse{
					ID:       item.ID,
					UserID:   item.User.ID,
					UserName: item.User.Name,
					Title:    item.Title,
					Address:  item.Address,
					City:     item.City,
					Price:    item.Price,
					Rating: rating,
					Features: featuresData,
					Ratings: ratingData,
				},
			)
		}

		return c.JSON(http.StatusOK, common.SuccessResponse(data))
	}
}

func (hc HouseController) GetHouseController() echo.HandlerFunc {

	return func(c echo.Context) error {

		id, _ := strconv.Atoi(c.Param("id"))

		house, err := hc.Repo.Get(id)

		if err != nil {
			return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
		}

		featuresData := []FeatureResponse{}
		for _, feature := range house.Features {
			featuresData = append(featuresData, FeatureResponse{
				ID:   feature.ID,
				Name: feature.Name,
			})
		}

		ratingData := []rating.RatingResponse{}
		ratings := []int{}

		for _, r := range house.Ratings {
			ratingData = append(ratingData, rating.RatingResponse{
				HouseID:  int(r.HouseID),
				UserID:   int(r.UserID),
				Username: r.User.Name,
				Rating:   r.Rating,
				Comment:  r.Comment,
			})

			ratings = append(ratings, r.Rating)
		}

		rating := helper.CalculateRatings(ratings)

		data := HouseResponse{
			ID:       house.ID,
			UserID:   house.User.ID,
			UserName: house.User.Name,
			Title:    house.Title,
			Address:  house.Address,
			City:     house.City,
			Price:    house.Price,
			Rating: rating,
			Features: featuresData,
			Ratings: ratingData,
		}

		return c.JSON(http.StatusOK, common.SuccessResponse(data))
	}
}

func (hc HouseController) UpdateHouseController() echo.HandlerFunc {

	return func(c echo.Context) error {
		putHouseReq := PutHouseRequestFormat{}
		id, _ := strconv.Atoi(c.Param("id"))
		user, _ := middleware.ExtractTokenUser(c)
		err := c.Bind(&putHouseReq)
		if err != nil {
			return err
		}

		newHouse := model.House{
			Title:   putHouseReq.Title,
			Address: putHouseReq.Address,
			City:    putHouseReq.City,
			Price:   putHouseReq.Price,
		}

		err = hc.Repo.HouseHasFeatureDelete(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		_, err = hc.Repo.Update(newHouse, id, user.UserID)
		if err != nil {
			return c.JSON(http.StatusNotFound, common.NewBadRequestResponse())
		}

		for _, feature := range putHouseReq.Features {
			err := hc.Repo.HouseHasFeature(model.HouseHasFeatures{
				HouseID:   newHouse.ID,
				FeatureID: uint(feature),
			})
			if err != nil {
				return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
			}
		}

		return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())
	}
}

func (hc HouseController) DeleteHouseCtrl() echo.HandlerFunc {

	return func(c echo.Context) error {
		var err error
		id, _ := strconv.Atoi(c.Param("id"))
		user, _ := middleware.ExtractTokenUser(c)

		err = hc.Repo.HouseHasFeatureDelete(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		_, err = hc.Repo.Delete(id, user.UserID)
		if err != nil {
			return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
		}

		return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())
	}
}
