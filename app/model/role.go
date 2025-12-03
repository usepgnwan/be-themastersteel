package model

type UserRole struct {
	ID        uint    `gorm:"primaryKey;autoIncrement;"json:"id" swaggerignore:"true"` // auto-increment
	Deskripsi *string `gorm:"size:255" json:"deskripsi"`
}

func (UserRole) TableName() string {
	return "tbl_role"
}
