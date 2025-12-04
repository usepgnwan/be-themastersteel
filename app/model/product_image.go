package model

import (
	"be-metalsteel/app/helpers"
	"time"

	"gorm.io/gorm"
)

type ProductImage struct {
	ID        string `gorm:"primaryKey" json:"id" form:"id" alias:"id" swaggerignore:"true"`
	ProductID string `gorm:"column:tbl_product_id;size:255;" json:"tbl_product_id" validate:"required"`
	FileName  string `json:"filename" column:"filename" alias:"filename" validate:"required"`
	Src       string `gorm:"column:src;type:text;" json:"src" validate:"required" example:"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAgAAAAIACAYAAAD0eNT6AAAACXBIWXM"`

	CreatedAt time.Time      `json:"created_at" swaggerignore:"true"`
	UpdatedAt time.Time      `json:"updated_at" swaggerignore:"true"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty" swaggerignore:"true"`
}

func (ProductImage) TableName() string {
	return "tbl_product_images"
}

func (d *ProductImage) BeforeCreate(tx *gorm.DB) (err error) {

	id := helpers.GenerateID()
	d.ID = id

	return nil
}
