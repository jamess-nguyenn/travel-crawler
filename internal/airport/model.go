package airport

import (
	"travel-crawler/internal/model"
)

const tableName = "airports"

type Airport struct {
	model.Model
	IATACode string `json:"iata_code" gorm:"column:iata_code;size:100;not null"`
	City     string `json:"city" gorm:"column:city;size:255;not null"`
	Country  string `json:"country" gorm:"column:country;size:255;not null"`
}

func (Airport) TableName() string {
	return tableName
}
