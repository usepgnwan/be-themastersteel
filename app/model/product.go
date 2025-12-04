package model

import (
	"be-metalsteel/app/helpers"
	"fmt"
	"time"

	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type Product struct {
	ID         string  `gorm:"primaryKey" json:"id" form:"id" alias:"id" swaggerignore:"true"`
	Title      string  `gorm:"column:title;size:255;" json:"title" example:"testing product" validate:"required"`
	Slug       string  `gorm:"column:slug;type:text;" json:"slug" swaggerignore:"true" `
	Descrption string  `gorm:"column:deskripsi;type:text;" json:"deskripsi"`
	Harga      string  `gorm:"column:harga;size:255;" json:"harga"`
	Category   *string `gorm:"column:category;" json:"category" validate:"required"`

	ProductImage []ProductImage `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE" json:"product_image" swaggerignore:"true"`

	CreatedAt time.Time      `json:"created_at" swaggerignore:"true"`
	UpdatedAt time.Time      `json:"updated_at" swaggerignore:"true"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty" swaggerignore:"true"`
}

func (Product) TableName() string {
	return "tbl_product"
}

func (d *Product) BeforeCreate(tx *gorm.DB) (err error) {
	baseSlug := slug.Make(d.Title)
	slugStr := baseSlug
	var count int64
	i := 1

	for {
		err = tx.Table(d.TableName()).Where("slug = ? AND deleted_at IS NULL", slugStr).Count(&count).Error
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		if count == 0 {
			break
		}

		slugStr = fmt.Sprintf("%s-%d", baseSlug, i)
		i++
	}

	id := helpers.GenerateID()
	d.ID = id
	d.Slug = slugStr
	return nil
}

func (d *Product) BeforeDelete(tx *gorm.DB) (err error) {
	return tx.Model(&ProductImage{}).Where("tbl_product_id = ?", d.ID).Delete(nil).Error
}
