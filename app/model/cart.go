package model

import (
	"be-metalsteel/app/helpers"
	"time"

	"gorm.io/gorm"
)

type Cart struct {
	ID           string  `gorm:"primaryKey" json:"id" form:"id" alias:"id" swaggerignore:"true"`
	UserID       string  `gorm:"column:user_id;size:255;" json:"user_id"`
	ProductID    string  `gorm:"column:product_id;size:255;" json:"product_id"`
	Title        string  `gorm:"column:title;size:255;" json:"title" example:"testing product" validate:"required"`
	Slug         string  `gorm:"column:slug;type:text;" json:"slug"`
	Descrption   string  `gorm:"column:deskripsi;type:text;" json:"deskripsi"`
	Status       string  `gorm:"column:status;type:text;" json:"status" swaggerignore:"true"`
	Harga        string  `gorm:"column:harga;size:255;" json:"harga"`
	Qty          string  `gorm:"column:qty;size:255;" json:"qty"`
	TotalHarga   string  `gorm:"column:total_harga;size:255;" json:"total_harga"`
	Category     *string `gorm:"column:category;" json:"category" validate:"required"`
	ProductImage string  `gorm:"column:product_image;type:text;" json:"product_image"`

	CreatedAt time.Time      `json:"created_at" swaggerignore:"true"`
	UpdatedAt time.Time      `json:"updated_at" swaggerignore:"true"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty" swaggerignore:"true"`
}

func (Cart) TableName() string {
	return "tbl_cart"
}

func (d *Cart) BeforeCreate(tx *gorm.DB) (err error) {
	id := helpers.GenerateID()
	d.ID = id
	d.Status = "cart"
	return nil
}
