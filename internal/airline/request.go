package airline

import "travel-crawler/internal/request"

type GetRequest struct {
	ID *uint `path:"id" json:"id"`
}

type ListRequest struct {
	request.ListRequest
	Alliance *string `query:"alliance" json:"alliance"`
}
