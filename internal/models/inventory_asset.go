package models

type InventoryAsset struct {
	Id         int    `gorm:"column:id;primary_key;not_null" json:"id" db:"id"`
	UserName   string `gorm:"not_null" json:"username" db:"username"`
	Name       string `gorm:"not_null" json:"name" db:"name"`
	SerialCode string `gorm:"not_null" json:"serial_code" db:"serial_code"`
	IsActive   bool   `gorm:"not_null" json:"is_active" db:"is_active"`
	Price      int    `gorm:"not_null" json:"price" db:"price"`
	Url        string `gorm:"not_null" json:"url" db:"url"`
}

// All inventory assets paginated response
type InventoryAssetList struct {
	Items      *[]InventoryAsset `json:"items"`
	TotalCount int64             `json:"total"`
	Page       int               `json:"page"`
	Size       int               `json:"size"`
	TotalPages int               `json:"pages"`
}
