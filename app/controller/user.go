package controller

import (
	. "be-metalsteel/app/helpers"
	"be-metalsteel/app/model"
	"be-metalsteel/connection"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// @Tags User
// @Summary list user
// @Description  list user
// @Param page query int false "(default : 1)"
// @Param limit query int false "(default : 10)"
// @Param name query string false "(optional)"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /api/user [get]
// @Param x-tes-themastersteel header string true "API secret key" default(UspGnwnelpsrSVTsQYu8LVRyGcl5m7kmi)
func GetUser(c echo.Context) error {

	data := &Paginate{
		Model: &model.User{},
	}
	db := connection.DB
	query := db.Model(&model.User{}).Preload("UserRole")

	name := c.QueryParam("name")

	if name != "" {
		query = query.Where("name ILIKE  ?", "%"+name+"%")
	}

	result := data.Paginate(query, c)

	return c.JSON(http.StatusOK, Response{Message: "success get data", Status: true, Data: result})
}

// @Tags User
// @Summary registrasi user
// @Accept json
// @Produce json
// @Accept multipart/form-data
// @Produce json
// @Param userrole body model.User true "add data"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /api/user [post]
// @Param x-tes-themastersteel header string true "API secret key" default(UspGnwnelpsrSVTsQYu8LVRyGcl5m7kmi)
func PostUser(c echo.Context) error {
	data := new(model.User)
	res := new(Response)
	res.Message = "Internal server error"
	res.Status = false

	if err := c.Bind(data); err != nil {
		return c.JSON(http.StatusInternalServerError, res)
	}

	valErr, err := ValidateData(data)
	if err != nil {
		res.Message = err.Error()
		return c.JSON(http.StatusInternalServerError, res)
	}

	if len(valErr) > 0 {
		res.Message = "Validation failed"
		res.Data = valErr
		return c.JSON(http.StatusBadRequest, res)
	}
	password := data.Password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Message: "Error while hashing password",
			Status:  false,
			Data:    nil,
		})
	}
	data.Password = string(hash)
	if errCreate := connection.DB.Create(&data).Error; errCreate != nil {
		res.Message = errCreate.Error()
		code := http.StatusInternalServerError
		return c.JSON(code, res)
	}

	if err := connection.DB.Model(&model.User{}).Preload("UserRole").Where("id = ?", data.ID).First(&data).Error; err != nil {
		fmt.Println("gagal get baru")
	}
	res.Status = true
	res.Message = "success create data"
	res.Data = data
	return c.JSON(http.StatusOK, res)
}

type DataLogin struct {
	UserContact string `json:"user_contact" form:"user_contact" alias:"user kontak" validate:"required,min=3"`
	Password    string `json:"password" form:"password" alias:"password" validate:"required,min=3"`
}

type JWTUSER struct {
	ID       string         ` json:"id" form:"id" alias:"id"`
	Name     *string        `json:"name" `
	Username *string        `json:"username" `
	Email    *string        `json:"email" `
	Phone    *string        `json:"phone" `
	RoleId   uint           `json:"role_id" `
	UserRole model.UserRole `json:"user_roles" gorm:"foreignKey:RoleId"`

	jwt.RegisteredClaims
}

type ResponseJWT struct {
	Message string      `json:"message" form:"message"`
	Status  bool        `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Token   *string     `json:"token"`
}

// @Tags User
// @Summary login ke user
// @Accept json
// @Produce json
// @Accept multipart/form-data
// @Produce json
// @Param data body DataLogin true "data login"
// @Success 200 {object} ResponseJWT
// @Failure 400 {object} ResponseJWT
// @Router /api/user/login  [post]
// @Param x-tes-themastersteel header string true "API secret key" default(UspGnwnelpsrSVTsQYu8LVRyGcl5m7kmi)
func UserLogin(e echo.Context) error {
	data := new(DataLogin)
	res := new(ResponseJWT)
	res.Data = nil
	res.Status = false
	res.Token = nil

	if err := e.Bind(data); err != nil {
		res.Message = err.Error()
		return e.JSON(http.StatusInternalServerError, res)
	}

	validationErr, err := ValidateData(data)

	if err != nil {
		res.Message = err.Error()
		return e.JSON(http.StatusInternalServerError, res)
	}

	if len(validationErr) > 0 {
		res.Message = "Validation failed"
		res.Data = validationErr
		return e.JSON(http.StatusBadRequest, res)
	}

	check := CheckUserName(data.UserContact)
	var usercontact = make(map[string]string)

	usercontact[check] = data.UserContact
	var user model.User
	if err := connection.DB.Model(&model.User{}).Where(usercontact).Preload("UserRole").First(&user).Where(usercontact).Error; err != nil {
		res.Message = err.Error()
		return e.JSON(http.StatusUnauthorized, res)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))
	if err != nil {
		res.Message = "validation Failed"
		res.Data = map[string]string{"user_contact": "email/nohp/username/password yang dimasukan salah"}
		return e.JSON(http.StatusUnauthorized, res)
	}

	claims := &JWTUSER{
		user.ID,
		user.Name,
		user.Username,
		user.Email,
		user.Phone,
		user.RoleId,
		user.UserRole,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72 * 30)), // expired on 30 days
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as ResponseJWT.
	t, err := token.SignedString([]byte(os.Getenv("SECRETKEY")))
	if err != nil {
		res.Message = err.Error()
		return e.JSON(http.StatusInternalServerError, res)
	}

	res.Message = "Login success"
	res.Status = true
	res.Data = map[string]interface{}{
		"role": user.UserRole,
	}
	res.Token = &t
	return e.JSON(http.StatusOK, res)
}
