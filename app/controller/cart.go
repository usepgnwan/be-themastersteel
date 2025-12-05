package controller

import (
	. "be-metalsteel/app/helpers"
	"be-metalsteel/app/model"
	"be-metalsteel/app/utils"
	"be-metalsteel/connection"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

// @Tags Cart
// @Summary list cart
// @Description  list cart
// @Param user_id query string true "(wajib)"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /api/cart [get]
// @Param x-tes-themastersteel header string true "API secret key" default(UspGnwnelpsrSVTsQYu8LVRyGcl5m7kmi)
func GetCart(c echo.Context) error {
	user_id := c.QueryParam("user_id")

	if user_id == "" {
		return c.JSON(http.StatusBadRequest, Response{
			Status:  false,
			Message: "user id wajib",
		})
	}
	UsingCrudHelper = &utils.Crud{
		Model: model.Cart{},
		Where: map[string]interface{}{"status": "cart", "user_id": user_id},
	}
	return UsingCrudHelper.Get(c)
}

// @Tags Cart
// @Summary add Cart
// @Accept json
// @Produce json
// @Accept multipart/form-data
// @Produce json
// @Param cart body model.Cart true "add data"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /api/cart [post]
// @Param x-tes-themastersteel header string true "API secret key" default(UspGnwnelpsrSVTsQYu8LVRyGcl5m7kmi)
func PostCart(c echo.Context) error {
	var d model.Cart
	body, _ := io.ReadAll(c.Request().Body)
	fmt.Println("RAW BODY:", string(body))
	c.Request().Body = io.NopCloser(bytes.NewBuffer(body))

	if err := json.Unmarshal(body, &d); err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Status:  false,
			Message: "JSON invalid: " + err.Error(),
		})
	}

	var count int64
	connection.DB.
		Model(&model.Cart{}).
		Where("status = 'cart'").
		Where("product_id = ?", d.ProductID).
		Where("user_id = ?", d.UserID).
		Count(&count)

	if count > 0 {
		return c.JSON(http.StatusBadRequest, Response{
			Status:  false,
			Message: "Product sudah ada di cart",
		})
	}

	UsingCrudHelper = &utils.Crud{
		Model: model.Cart{},
	}
	id := GenerateID()
	column := "ID"
	return UsingCrudHelper.Create(c, &column, &id)
}
