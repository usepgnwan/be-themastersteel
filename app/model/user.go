package model

import (
	helpers "be-metalsteel/app/helpers"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID       string   `gorm:"primaryKey" json:"id" form:"id" alias:"id" swaggerignore:"true"`
	Name     *string  `gorm:"size:255" json:"name" validate:"required"`
	Username *string  `gorm:"size:255" json:"username"`
	Email    *string  `gorm:"size:255" json:"email"`
	Phone    *string  `gorm:"size:255" json:"phone"`
	Password string   `gorm:"column:password" json:"password"`
	RoleId   uint     `gorm:"size:255" json:"role_id" validate:"required"`
	UserRole UserRole `json:"user_roles" gorm:"foreignKey:RoleId" swaggerignore:"true"`

	// gorm.Model
	CreatedAt time.Time      `json:"created_at" swaggerignore:"true"`
	UpdatedAt time.Time      `json:"updated_at" swaggerignore:"true"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty" swaggerignore:"true"`
}

func (User) TableName() string {
	return "tbl_users"
}

func (d *User) BeforeCreate(tx *gorm.DB) (err error) {
	d.ID = helpers.GenerateID()
	return
}
