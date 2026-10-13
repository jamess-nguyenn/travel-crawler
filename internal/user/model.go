package user

import (
	"time"
	"travel-crawler/internal/model"
)

const tableName = "users"

type Status string

type User struct {
	model.Model
	Email           string     `json:"email" gorm:"column:email;size:255;not null;uniqueIndex"`
	Auth0ID         string     `json:"auth0_id" gorm:"column:auth0_id;size:255;uniqueIndex"`
	EmailVerifiedAt *time.Time `json:"email_verified_at" gorm:"column:email_verified_at"`
	PasswordHash    *string    `json:"-" gorm:"column:password_hash;size:255"`
	FirstName       string     `json:"first_name" gorm:"column:first_name;size:100;not null"`
	LastName        string     `json:"last_name" gorm:"column:last_name;size:100;not null"`
	FullName        string     `json:"full_name" gorm:"column:full_name;size:255;not null"`
	PhoneNumber     *string    `json:"phone_number" gorm:"column:phone_number;size:100"`
	Status          Status     `json:"status" gorm:"column:status;size:100;not null;default:active"`
	AvatarURL       *string    `json:"avatar_url" gorm:"column:avatar_url;size:500"`
	LastLoginAt     *time.Time `json:"last_login_at" gorm:"column:last_login_at"`
}

func (u User) TableName() string {
	return tableName
}

func BuildFullName(firstName, lastName string) string {
	return firstName + " " + lastName
}

func (u User) GetFullName() string {
	return BuildFullName(u.FirstName, u.LastName)
}
