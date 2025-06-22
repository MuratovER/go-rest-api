package models

// InventoryAsset represents an inventory asset entity
type InventoryAsset struct {
	ID         int    `gorm:"column:id;primary_key;not_null" json:"id" db:"id"`
	UserName   string `gorm:"not_null" json:"username" db:"username"`
	Name       string `gorm:"not_null" json:"name" db:"name"`
	SerialCode string `gorm:"not_null" json:"serial_code" db:"serial_code"`
	IsActive   bool   `gorm:"not_null" json:"is_active" db:"is_active"`
	Price      int    `gorm:"not_null" json:"price" db:"price"`
	URL        string `gorm:"not_null" json:"url" db:"url"`
}

// TableName specifies the table name for InventoryAsset
func (InventoryAsset) TableName() string {
	return "inventory_assets"
}

// InventoryAssetList represents a paginated response of inventory assets
type InventoryAssetList struct {
	Items      *[]InventoryAsset `json:"items"`
	TotalCount int64             `json:"total"`
	Page       int               `json:"page"`
	Size       int               `json:"size"`
	TotalPages int               `json:"pages"`
}
