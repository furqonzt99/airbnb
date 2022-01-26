package rating

import (
	"net/http"
	"strconv"

	"github.com/furqonzt99/airbnb/delivery/common"
	mw "github.com/furqonzt99/airbnb/delivery/middleware"
	"github.com/furqonzt99/airbnb/model"
	rr "github.com/furqonzt99/airbnb/repository/rating"
	"github.com/labstack/echo/v4"
)

type RatingController struct {
	Repository rr.RatingRepository
}

func NewRatingController(repo rr.RatingRepository) *RatingController {
	return &RatingController{Repository: repo}
}

func (rc RatingController) Create(c echo.Context) error {
	var ratingRequest PostRatingRequest

	if err := c.Bind(&ratingRequest); err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	if err := c.Validate(&ratingRequest); err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	user, _ := mw.ExtractTokenUser(c)

	data := model.Rating{
		HouseID: uint(ratingRequest.HouseID),
		UserID:  uint(user.UserID),
		Rating:  ratingRequest.Rating,
		Comment: ratingRequest.Comment,
	}

	ratingData, err := rc.Repository.Create(data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	return c.JSON(http.StatusOK, common.SuccessResponse(ratingData)) 
}

func (rc RatingController) Update(c echo.Context) error {
	var ratingRequest UpdateRatingRequest

	if err := c.Bind(&ratingRequest); err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	if err := c.Validate(&ratingRequest); err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	houseId, err := strconv.Atoi(c.Param("houseId"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	user, _ := mw.ExtractTokenUser(c)

	data := model.Rating{
		HouseID: uint(houseId),
		UserID:  uint(user.UserID),
		Rating:  ratingRequest.Rating,
		Comment: ratingRequest.Comment,
	}

	ratingData, err := rc.Repository.Update(data)
	if err != nil {
		return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
	}

	return c.JSON(http.StatusOK, common.SuccessResponse(ratingData)) 
}

func (rc RatingController) Delete(c echo.Context) error {

	houseId, err := strconv.Atoi(c.Param("houseId"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	user, _ := mw.ExtractTokenUser(c)

	ratingData, err := rc.Repository.Delete(user.UserID, houseId)
	if err != nil {
		return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
	}

	return c.JSON(http.StatusOK, common.SuccessResponse(ratingData)) 
}