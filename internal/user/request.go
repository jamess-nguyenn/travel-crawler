package user

import "travel-crawler/internal/request"

type UpdateRequest struct {
	FirstName   string  `json:"first_name" validate:"required,max=100"`
	LastName    string  `json:"last_name" validate:"required,max=100"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty,max=100"`
	Status      Status  `json:"status" validate:"required,oneof=active inactive suspended"`
}

type GetRequest struct {
	ID *uint `path:"id" json:"id"`
}

type ListAllRequest struct {
	IDs    []uint64 `query:"ids" json:"ids"`
	Emails []string `query:"emails" json:"emails"`
	Status *Status  `query:"status" json:"status"`
}

type ListRequest struct {
	request.ListRequest
	Status *Status `query:"status" json:"status"`
}

type ExportRequest struct {
	request.ExportRequest
	Status *Status `query:"status" json:"status"`
}
