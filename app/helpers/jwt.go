package helpers

import "github.com/golang-jwt/jwt/v5"

// jwtCustomClaims are custom claims extending default ones.
// See https://github.com/golang-jwt/jwt for more examples

type MyJwt struct {
	ID          string  `json:"id"`
	RoleId      *string `json:"role_id"`
	Name        *string `json:"name"`
	ProfileName *string `json:"profile_name"`
	Email       *string `json:"email"`
	Phone       *string `json:"phone"`
	AppsGroup   *string `json:"apps_group"`
	jwt.RegisteredClaims
}
