package airline

import (
	"travel-crawler/internal/model"
)

const tableName = "airlines"

type Airline struct {
	model.Model
	ProviderCode string  `json:"provider_code" gorm:"column:provider_code;size:100;not null"`
	IATACode     string  `json:"iata_code" gorm:"column:iata_code;size:100;not null"`
	Name         string  `json:"name" gorm:"column:name;size:255;not null"`
	Alliance     *string `json:"alliance" gorm:"column:alliance;size:100"`
	Active       bool    `json:"active" gorm:"column:active;not null;default:true"`
	Logo         *string `json:"logo" gorm:"column:logo;size:500"`
}

func (Airline) TableName() string {
	return tableName
}
