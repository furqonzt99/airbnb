package feature

import (
	"net/http"

	"github.com/furqonzt99/airbnb/delivery/common"
	"github.com/furqonzt99/airbnb/repository/feature"
	"github.com/labstack/echo/v4"
)

type FeatureController struct {
	Repo feature.FeatureInterface
}

func NewFeatureControllers(fr feature.FeatureInterface) *FeatureController {
	return &FeatureController{Repo: fr}
}

func (fc FeatureController) GetAllFeatureCtrl() echo.HandlerFunc {

	return func(c echo.Context) error {

		features, _ := fc.Repo.GetAll()

		data := []FeatureResponse{}
		for _, item := range features {
			data = append(
				data, FeatureResponse{
					ID:   item.ID,
					Name: item.Name,
				},
			)
		}

		return c.JSON(
			http.StatusOK, common.SuccessResponse(data),
		)
	}

}
