package response

type Pagination struct {
	Page      int   `json:"page"`
	PageSize  int   `json:"page_size"`
	Total     int64 `json:"total"`
	TotalPage int   `json:"total_page"`
}

func NewPagination(page, pageSize int, total int64) Pagination {
	totalPage := 0

	if pageSize > 0 {
		totalPage = int((total + int64(pageSize) - 1) / int64(pageSize))
	}

	return Pagination{
		Page:      page,
		PageSize:  pageSize,
		Total:     total,
		TotalPage: totalPage,
	}
}
