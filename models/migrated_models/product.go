package migrated_models

import (
	"github.com/google/uuid"
	"time"
)

type Product struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	Price       int       `json:"price"`
	ImageURL    string    `json:"image_url"`                     // URL to the skin/image
	IsActive    bool      `gorm:"default:true" json:"is_active"` // To enable/disable products
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type Purchase struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid" json:"user_id"`
	ProductID uuid.UUID `gorm:"type:uuid" json:"product_id"`
	PricePaid int       `json:"price_paid"` // Points spent (in case price changes later)
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	User    User    `gorm:"foreignKey:UserID" json:"user"`
	Product Product `gorm:"foreignKey:ProductID" json:"product"`
}

type UserInventory struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid" json:"user_id"`
	ProductID uuid.UUID `gorm:"type:uuid" json:"product_id"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`

	User    User    `gorm:"foreignKey:UserID" json:"user"`
	Product Product `gorm:"foreignKey:ProductID" json:"product"`
}
