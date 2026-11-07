package response

type ListResponse[T any] struct {
	Items      []T         `json:"items"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Request    any         `json:"request"`
}
