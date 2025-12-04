package controller

import (
	. "be-metalsteel/app/helpers"
	"be-metalsteel/app/model"
	"be-metalsteel/app/utils"
	"be-metalsteel/connection"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

// @Tags Product
// @Summary list product
// @Description  list product
// @Param page query int false "(default : 1)"
// @Param limit query int false "(default : 10)"
// @Param title query string false "(optional)"
// @Param category query []string false "Filter kategori (bisa multiple)"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /api/product [get]
// @Param x-tes-themastersteel header string true "API secret key" default(UspGnwnelpsrSVTsQYu8LVRyGcl5m7kmi)
func GetProduct(c echo.Context) error {

	data := &Paginate{
		Model: &model.Product{},
	}
	db := connection.DB
	query := db.Model(&model.Product{}).Preload("ProductImage")

	title := c.QueryParam("title")

	if title != "" {
		query = query.Where("title ILIKE  ?", "%"+title+"%")
	}

	categoriesRaw := c.QueryParams()["category"]

	var categories []string

	if len(categoriesRaw) > 0 {
		categories = strings.Split(categoriesRaw[0], ",")
	}

	if len(categories) > 0 {
		query = query.Where("category IN ?", categories)
	}
	result := data.Paginate(query, c)

	return c.JSON(http.StatusOK, Response{Message: "success get data", Status: true, Data: result})
}

var UsingCrudHelper *utils.Crud

// @Tags Product
// @Summary add product
// @Accept json
// @Produce json
// @Accept multipart/form-data
// @Produce json
// @Param product body model.Product true "add data"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /api/product [post]
// @Param x-tes-themastersteel header string true "API secret key" default(UspGnwnelpsrSVTsQYu8LVRyGcl5m7kmi)
func PostProduct(c echo.Context) error {
	UsingCrudHelper = &utils.Crud{
		Model: model.Product{},
	}
	id := GenerateID()
	column := "ID" // sesuaikan degan struct
	return UsingCrudHelper.Create(c, &column, &id)
}

// @Tags Product
// @Summary login ke user
// @Accept json
// @Produce json
// @Accept multipart/form-data
// @Produce json
// @Param data body model.ProductImage true "data image"
// @Success 200 {object} ResponseJWT
// @Failure 400 {object} ResponseJWT
// @Router /api/product-images [post]
// @Param x-tes-themastersteel header string true "API secret key" default(UspGnwnelpsrSVTsQYu8LVRyGcl5m7kmi)
func PostProductImage(c echo.Context) error {

	var data model.ProductImage
	if err := c.Bind(&data); err != nil {
		return c.JSON(http.StatusBadRequest, Response{Status: false, Message: "rererere"})
	}
	valErr, err := ValidateData(data)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Message: "failed", Status: false})
	}
	if len(valErr) > 0 {
		return c.JSON(http.StatusBadRequest, Response{Message: "validation failed", Status: false})
	}
	// files, err := GetDataFile64(data.Src)
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, Response{Status: false, Message: err.Error()})
	// }

	tmpThumb, err := Base64ToFile(data.Src, data.FileName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Message: "Error decoding base64 thumbnail: " + err.Error(),
			Status:  false,
		})
	}

	defer os.Remove(tmpThumb.Name())

	nameThumb := data.FileName

	uploadDir := "uploads/product-images"
	os.MkdirAll(uploadDir, 0755)

	dstPath := fmt.Sprintf("%s/%d-%s", uploadDir, time.Now().Unix(), nameThumb)

	// buka file temporary
	srcFile, err := os.Open(tmpThumb.Name())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Message: err.Error(),
			Status:  false,
		})
	}
	defer srcFile.Close()

	// buat file destination
	dst, err := os.Create(dstPath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Message: err.Error(),
			Status:  false,
		})
	}
	defer dst.Close()

	// salin isi file
	if _, err := io.Copy(dst, srcFile); err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Message: err.Error(),
			Status:  false,
		})
	}

	// URL final yg bisa diakses client
	fileURL := os.Getenv("BASE_URL") + dstPath

	data.Src = fileURL
	if err := connection.DB.Model(&model.ProductImage{}).Create(&data).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Message: err.Error(), Status: false})
	}

	return c.JSON(http.StatusOK, Response{
		Message: "Image berhasil diupload",
		Data:    data,
		Status:  true,
	})
}

// @Tags Product
// @Summary get detail product
// @Accept json
// @Produce json
// @Param slug path string true "(wajib)"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /api/product/{slug} [get]
// @Param x-tes-themastersteel header string true "API secret key" default(UspGnwnelpsrSVTsQYu8LVRyGcl5m7kmi)
func GetIdProduct(c echo.Context) error {
	slug := c.Param("slug")
	if slug == "" {
		return c.JSON(http.StatusBadGateway, Response{Message: "slug wajib di isi", Status: false})
	}

	var result model.Product
	db := connection.DB.Model(&model.Product{}).Preload("ProductImage")

	if err := db.Where("slug =?", slug).First(&result).Error; err != nil {
		return c.JSON(http.StatusNotFound, Response{Message: err.Error(), Status: false})
	}

	return c.JSON(http.StatusOK, Response{Message: "success get data", Status: true, Data: result})
}
