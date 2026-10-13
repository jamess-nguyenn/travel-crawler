package model

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type Model struct {
	ID        uint64         `json:"id" gorm:"column:id;primaryKey"`
	Metadata  datatypes.JSON `json:"metadata" gorm:"column:metadata;type:json"`
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at;not null"`
	CreatedBy uint64         `json:"created_by" gorm:"column:created_by;not null"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at;not null"`
	UpdatedBy uint64         `json:"updated_by" gorm:"column:updated_by;not null"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at"`
	DeletedBy *uint64        `json:"deleted_by" gorm:"column:deleted_by"`
}

type TimestampModel struct {
	ID        uint64    `json:"id" gorm:"column:id;primaryKey"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;not null"`
}
