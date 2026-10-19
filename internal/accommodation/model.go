package accommodation

import (
	"travel-crawler/internal/model"
)

const tableName = "accommodations"

type Accommodation struct {
	model.Model
	ProviderCode     string   `json:"provider_code" gorm:"column:provider_code;size:100;not null"`
	ExternalID       string   `json:"external_id" gorm:"column:external_id;size:255;not null"`
	ExternalParentID *string  `json:"external_parent_id" gorm:"column:external_parent_id;size:255"`
	Name             string   `json:"name" gorm:"column:name;size:255;not null"`
	Description      string   `json:"description" gorm:"column:description;type:text;not null"`
	PropertyType     string   `json:"property_type" gorm:"column:property_type;size:100;not null"`
	Chain            string   `json:"chain" gorm:"column:chain;size:255;not null"`
	Currency         string   `json:"currency" gorm:"column:currency;type:char(3);not null"`
	Country          string   `json:"country" gorm:"column:country;type:char(2);not null"`
	City             string   `json:"city" gorm:"column:city;size:255;not null"`
	Latitude         *float64 `json:"latitude" gorm:"column:latitude;type:decimal(10,7)"`
	Longitude        *float64 `json:"longitude" gorm:"column:longitude;type:decimal(10,7)"`
	Address          string   `json:"address" gorm:"column:address;size:500;not null"`
	Zip              string   `json:"zip" gorm:"column:zip;size:50;not null"`
	MainPhoto        *string  `json:"main_photo" gorm:"column:main_photo;size:500"`
	Thumbnail        *string  `json:"thumbnail" gorm:"column:thumbnail;size:500"`
	Stars            float64  `json:"stars" gorm:"column:stars;not null"`
	Rating           *float64 `json:"rating" gorm:"column:rating;type:decimal(3,1)"`
	ReviewCount      uint32   `json:"review_count" gorm:"column:review_count;not null;default:0"`
}

func (Accommodation) TableName() string {
	return tableName
}
