package controller

import (
	. "be-metalsteel/app/helpers"
	"be-metalsteel/app/model"
	"be-metalsteel/connection"
	"net/http"

	"github.com/labstack/echo/v4"
)

// @Tags Role
// @Summary list role
// @Description  list role
// @Param page query int false "(default : 1)"
// @Param limit query int false "(default : 10)"
// @Param deskripsi query string false "(optional)"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /api/role [get]
// @Param x-tes-themastersteel header string true "API secret key" default(UspGnwnelpsrSVTsQYu8LVRyGcl5m7kmi)
func GetRole(c echo.Context) error {

	data := &Paginate{
		Model: &model.UserRole{},
	}
	db := connection.DB
	query := db.Model(&model.UserRole{})

	deskripsi := c.QueryParam("deskripsi")

	if deskripsi != "" {
		query = query.Where("deskripsi ILIKE  ?", "%"+deskripsi+"%")
	}

	result := data.Paginate(query, c)

	return c.JSON(http.StatusOK, Response{Message: "success get data", Status: true, Data: result})
}
