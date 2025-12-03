package utils

import (
	"be-metalsteel/app/helpers"
	"be-metalsteel/connection"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-playground/validator"
	"gorm.io/gorm"

	"github.com/labstack/echo/v4"
)

type Crud struct {
	// DB    *gorm.DB
	Model        interface{}
	SelectFields []string
	Where        map[string]interface{}
	Option       string
	OnlyParam    map[string]bool
}
type Response struct {
	Message string      `json:"message" form:"message"`
	Status  bool        `json:"status"`
	Data    interface{} `json:"data"`
}
type ErrResponse struct {
	Message string `json:"message" form:"message"`
	Status  bool   `json:"status"`
}

var validate = validator.New()

type SwaggerMetadata struct {
	Summary     string
	Description string
	Tags        string
	Model       string
	Route       string
	Method      string
}

// @Summary create a history user agent
// @Description create a history user agent with the input payload
// @Tags SwaggerMetadata.Tags
// @Accept json
// @Produce json
// @Param user body model.Employee true "Employee details"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /api/user_agent  [get]
// @Param x-poolapack-header header string true "API secret key"

// Dynamic metadata generation
func (c *Crud) GenerateSwaggerMetadata(metadata SwaggerMetadata) string {
	fmt.Printf(`
		// @Summary      %s
		// @Description  %s
		// @Tags         %s
		// @Produce      json
		// @Success      200  {array}  %s "List of %s"
		// @Failure      400  {object}  Error "Bad request"
		// @Router       %s [%s]
	`, metadata.Summary, metadata.Description, metadata.Tags, "model.Employee", metadata.Tags, metadata.Route, metadata.Method)
	return fmt.Sprintf(`
	// @Summary      %s
	// @Description  %s
	// @Tags         %s
	// @Produce      json
	// @Success      200  {array}  %s "List of %s"
	// @Failure      400  {object}  Error "Bad request"
	// @Router       %s [%s]
`, metadata.Summary, metadata.Description, metadata.Tags, "model.Employee", metadata.Tags, metadata.Route, metadata.Method)

}

// Middleware to add dynamic Swagger metadata
func (c *Crud) MetadataMiddleware(metadata SwaggerMetadata) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(e echo.Context) error {
			// Generate metadata dynamically
			swaggerMetadata := c.GenerateSwaggerMetadata(metadata)
			fmt.Println(swaggerMetadata) // Print metadata or log it somewhere for documentation

			// Continue with the next handler
			return next(e)
		}
	}
}
func (c *Crud) Get(e echo.Context) error {
	id := e.Param("id")

	var data interface{}
	query := connection.DB.Unscoped().Model(c.Model)
	if len(c.SelectFields) > 0 {
		selected := strings.Join(c.SelectFields, ",")
		query = query.Select(selected)
	}

	// query where dinamic api
	queryString := e.QueryString()

	var qs = make(map[string]interface{})
	if len(c.OnlyParam) > 0 {
		if queryString != "" {
			qrSt := strings.Split(queryString, "&")
			for _, v := range qrSt {
				parts := strings.Split(v, "=")
				key := parts[0]
				value := parts[1]
				// Add the key-value pair to the map
				if _, isAllowed := c.OnlyParam[key]; isAllowed {
					qs[key] = value
				} else {
					return e.JSON(http.StatusBadGateway, ErrResponse{"Invalid query parameter: " + key, false})
				}
			}
		}

		if qs != nil {
			query.Where(qs)
		}
	}

	// Where condition hardcode
	if len(c.Where) > 0 {
		for key, value := range c.Where {
			// Dynamically add where clauses to the query hardcode
			query = query.Where(key, value)
		}
	}

	if id == "" {
		var result []map[string]interface{}
		// dinamyc query
		if err := query.Find(&result).Error; err != nil {
			return e.JSON(http.StatusInternalServerError, ErrResponse{"Failed to execute " + err.Error(), false})
		}
		if qs != nil {
			if result == nil {
				return e.JSON(http.StatusOK, Response{"Data not found", false, result})
			}
		}
		data = result
	} else {

		var result map[string]interface{}
		if err := query.Where("id = ?", id).First(&result).Error; err != nil {
			return e.JSON(http.StatusInternalServerError, ErrResponse{"Failed to find data " + err.Error(), false})
		}
		data = result
	}

	return e.JSON(http.StatusOK, Response{"success get data", true, data})
}
func (c *Crud) Create(e echo.Context, column *string, value *string) error {

	modelType := reflect.TypeOf(c.Model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem() // Get the underlying type if it is a pointer
	}

	// Create a new instance of the model (pointer type)
	data := reflect.New(modelType).Interface()

	if err := e.Bind(data); err != nil {
		fmt.Println(ErrResponse{err.Error(), false})
		return e.JSON(http.StatusInternalServerError, ErrResponse{err.Error(), false})
	}

	validationErrors, err := helpers.ValidateData(data)

	if err != nil {
		return e.JSON(http.StatusInternalServerError, ErrResponse{err.Error(), false})
	}

	if len(validationErrors) > 0 {
		return e.JSON(http.StatusInternalServerError, Response{Message: "Validation failed", Status: false, Data: validationErrors})
	}

	if column != nil && value != nil {
		v := reflect.ValueOf(data).Elem()
		field := v.FieldByName(*column)

		if field.IsValid() && field.CanSet() {
			switch field.Kind() {
			case reflect.String:
				field.SetString(*value)
			case reflect.Int, reflect.Int64:
				if intVal, err := strconv.ParseInt(*value, 10, 64); err == nil {
					field.SetInt(intVal)
				}
			case reflect.Uint, reflect.Uint64:
				if uintVal, err := strconv.ParseUint(*value, 10, 64); err == nil {
					field.SetUint(uintVal)
				}
			// Tambah tipe lain sesuai kebutuhan
			default:
				fmt.Println("Unsupported field type")
			}
		} else {
			fmt.Println("Field", *column, "not found or not settable")
		}
	}
	// Save the model to the database
	if err := connection.DB.Create(data).Error; err != nil {
		return e.JSON(http.StatusInternalServerError, ErrResponse{err.Error(), false})
	}

	return e.JSON(http.StatusOK, Response{"success adding record", true, data})
}

func (c *Crud) Update(e echo.Context) error {
	id := e.Param("id")
	if id == "" {
		return e.JSON(http.StatusBadRequest, Response{"Id param required", false, nil})
	}
	var input map[string]interface{}
	if err := json.NewDecoder(e.Request().Body).Decode(&input); err != nil {
		return e.JSON(http.StatusBadRequest, Response{"Invalid JSON format", false, nil})
	}

	if len(c.SelectFields) > 0 {
		checkinput := make(map[string]string)

		var contains = func(slice []string, item string) bool {
			for _, v := range slice {
				if v == item {
					return true
				}
			}
			return false
		}

		for key := range input {
			if !contains(c.SelectFields, key) {
				checkinput[key] = "Invalid key : " + key
			}
		}

		if len(checkinput) > 0 {
			return e.JSON(http.StatusBadRequest, Response{"Validation failed", false, checkinput})
		}
	}
	// Dynamically create a new instance of the model type
	modelType := reflect.TypeOf(c.Model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem() // Get the underlying type if it is a pointer
	}

	// Create a new instance of the model (pointer type) for binding the request data
	post := reflect.New(modelType).Interface()

	// Bind the incoming request body to the model
	bodyBytes, _ := json.Marshal(input) // encode ulang input map ke json
	if err := json.Unmarshal(bodyBytes, &post); err != nil {
		return e.JSON(http.StatusBadRequest, Response{"Invalid data", false, nil})
	}

	// Validate the model
	// if err := validate.Struct(post); err != nil {
	// 	validationErrors := make(map[string]string)
	// 	for _, e := range err.(validator.ValidationErrors) {
	// 		errorMessage := fmt.Sprintf("Key: '%s' Error:Field validation for '%s' failed on the '%s' tag", e.Field(), e.StructField(), e.Tag())
	// 		validationErrors[e.Field()] = errorMessage
	// 	}
	// 	return e.JSON(http.StatusBadRequest, Response{"Validation failed", false, validationErrors})
	// }

	validationErrors, err := helpers.ValidateData(post)

	if err != nil {
		return e.JSON(http.StatusInternalServerError, ErrResponse{err.Error(), false})
	}

	if len(validationErrors) > 0 {
		return e.JSON(http.StatusInternalServerError, Response{Message: "Validation failed", Status: false, Data: validationErrors})
	}

	// Dynamically create a new instance of the model for fetching existing data
	existing := reflect.New(modelType).Interface()

	// Fetch the existing record from the database
	if err := connection.DB.Unscoped().Model(c.Model).Where("id =?", id).First(existing).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return e.JSON(http.StatusNotFound, Response{
				"Data not found", false, nil,
			})
		}
		return e.JSON(http.StatusInternalServerError, Response{"Failed to fetch data: " + err.Error(), false, nil})
	}
	// fmt.Println(post)
	// Update the existing record with the new data
	if err := connection.DB.Model(existing).Updates(post).Error; err != nil {
		return e.JSON(http.StatusInternalServerError, Response{"Failed to update data: " + err.Error(), false, existing})
	}

	result := make(map[string]interface{})
	query2 := connection.DB.Unscoped().Model(c.Model).Where("id =?", id)

	if len(c.SelectFields) > 0 {
		selected := strings.Join(c.SelectFields, ",")
		query2 = query2.Select(selected)
	}

	if err := query2.First(result).Error; err != nil {
		return e.JSON(http.StatusInternalServerError, Response{"Failed to retrieve updated data: " + err.Error(), false, nil})
	}

	return e.JSON(http.StatusOK, Response{"Successfully updated data", true, result})
}

func (c *Crud) Delete(e echo.Context, column *string, status *string) error {
	id := e.Param("id")
	if id == "" {
		return e.JSON(http.StatusBadRequest, Response{"Id param required", false, nil})
	}

	col := "id"
	if column != nil {
		col = *column
	}
	// Dynamically create a new instance of the model type
	modelType := reflect.TypeOf(c.Model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem() // Get the underlying type if it is a pointer
	}

	// Dynamically create a new instance of the model for fetching existing data
	modelInstance := reflect.New(modelType).Interface()

	// Cek apakah data dengan ID tersebut ada
	if err := connection.DB.Unscoped().Model(c.Model).Where(col+"= ?", id).First(modelInstance).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return e.JSON(http.StatusNotFound, Response{"Data not found", false, nil})
		}
		return e.JSON(http.StatusInternalServerError, Response{"Failed to fetch data: " + err.Error(), false, nil})
	}
	validstat := "hard"

	if status != nil && *status == validstat {

		if err := connection.DB.Unscoped().Where(col+" = ?", id).Delete(&c.Model).Error; err != nil {
			return e.JSON(http.StatusInternalServerError, Response{"Failed to delete data: " + err.Error(), false, nil})
		}
	} else {
		fmt.Print("sss")
		if err := connection.DB.Where(col+" = ?", id).Delete(&c.Model).Error; err != nil {
			return e.JSON(http.StatusInternalServerError, Response{"Failed to delete data: " + err.Error(), false, nil})
		}
	}

	return e.JSON(http.StatusOK, Response{"Successfully delete data", true, nil})
}
