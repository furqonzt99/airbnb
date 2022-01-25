package house

import (
	"net/http"
	"strconv"

	"github.com/furqonzt99/airbnb/delivery/common"
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
		newProductReq := CreateHouseRequestFormat{}
		c.Bind(&newProductReq)

		newProduct := model.House{
			Title:    newProductReq.Title,
			Address:  newProductReq.Address,
			City:     newProductReq.City,
			Price:    newProductReq.Price,
			Features: newProductReq.Features,
		}

		product, err := hc.Repo.Create(newProduct)

		if err != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		return c.JSON(http.StatusOK, common.SuccessResponse(product))
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
			data = append(
				data, HouseResponse{
					ID:       item.ID,
					UserID:   item.User.ID,
					UserName: item.User.Name,
					Title:    item.Title,
					Address:  item.Address,
					City:     item.City,
					Price:    item.Price,
					Features: item.Features,
				},
			)
		}
		return c.JSON(http.StatusOK, common.PaginationResponse(page, perpage, data))
	}
}

func (hc HouseController) GetHouseController() echo.HandlerFunc {

	return func(c echo.Context) error {

		id, _ := strconv.Atoi(c.Param("id"))

		house, err := hc.Repo.Get(id)

		if err != nil {
			return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
		}

		data := HouseResponse{
			ID:       house.ID,
			UserID:   house.User.ID,
			UserName: house.User.Name,
			Title:    house.Title,
			Address:  house.Address,
			City:     house.City,
			Price:    house.Price,
			Features: house.Features,
		}

		return c.JSON(http.StatusOK, common.SuccessResponse(data))
	}
}

func (hc HouseController) UpdateHouseController() echo.HandlerFunc {

	return func(c echo.Context) error {
		PutHouseReq := PutHouseRequestFormat{}
		id, _ := strconv.Atoi(c.Param("id"))
		err := c.Bind(&PutHouseReq)
		if err != nil {
			return err
		}

		newHouse := model.House{
			Title:    PutHouseReq.Title,
			Address:  PutHouseReq.Address,
			City:     PutHouseReq.City,
			Price:    PutHouseReq.Price,
			Features: PutHouseReq.Features,
		}

		result, err := hc.Repo.Update(newHouse, id)
		if err != nil {
			return c.JSON(http.StatusNotFound, common.NewBadRequestResponse())
		}

		data := HouseResponse{
			ID:       result.ID,
			UserID:   result.UserID,
			UserName: result.User.Name,
			Title:    result.Title,
			Address:  result.Address,
			City:     result.City,
			Price:    result.Price,
			Features: result.Features,
		}

		return c.JSON(http.StatusOK, common.SuccessResponse(data))
	}
}

func (hc HouseController) DeleteHouseCtrl() echo.HandlerFunc {

	return func(c echo.Context) error {
		var err error
		id, _ := strconv.Atoi(c.Param("id"))

		_, err = hc.Repo.Delete(id)
		if err != nil {
			return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
		}

		return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())
	}
}
