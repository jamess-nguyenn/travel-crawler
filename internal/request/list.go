package request

import (
	"strings"
	"time"
	"travel-crawler/internal/helper"
)

const (
	DefaultPage     = 1
	DefaultPageSize = 10

	PageSizeAll = 10000

	DefaultSortBy = "id"
	SortOrderAsc  = "asc"
	SortOrderDesc = "desc"

	ExportFormatCSV     = "csv"
	ExportFormatExcel   = "xlsx"
	DefaultExportFormat = ExportFormatCSV
)

type SortRequest interface {
	GetSortBy() string
	GetSortOrder() string
}

type PaginationRequest interface {
	GetPage() int
	GetPageSize() int
	NoLimit() bool
}

type QueryOptions struct {
	Search    *string    `query:"search" json:"search"`
	DateFrom  *time.Time `query:"date_from" json:"date_from"`
	DateTo    *time.Time `query:"date_to" json:"date_to"`
	SortBy    *string    `query:"sort_by" json:"sort_by"`
	SortOrder *string    `query:"sort_order" json:"sort_order"`
}

func (q *QueryOptions) NormalizeDateRange() {
	if q.DateTo != nil {
		dateTo := helper.EndOfDay(*q.DateTo)
		q.DateTo = &dateTo
	}
}

func (q QueryOptions) GetSortBy() string {
	value := helper.StringValue(q.SortBy)
	if value == "" {
		return DefaultSortBy
	}

	return value
}

func (q QueryOptions) GetSortOrder() string {
	value := strings.ToLower(helper.StringValue(q.SortOrder))

	if value == SortOrderAsc {
		return SortOrderAsc
	}

	return SortOrderDesc
}

type ListRequest struct {
	QueryOptions
	Page     *int `query:"page" json:"page"`
	PageSize *int `query:"page_size" json:"page_size"`
}

func (r ListRequest) GetPage() int {
	value := helper.IntValue(r.Page)
	if value == 0 {
		return DefaultPage
	}

	return value
}

func (r ListRequest) GetPageSize() int {
	value := helper.IntValue(r.PageSize)
	if value == 0 {
		return DefaultPageSize
	}

	return value
}

func (r ListRequest) NoLimit() bool {
	return r.GetPageSize() >= PageSizeAll
}

func (r ListRequest) GetOffset() int {
	return helper.CalculateOffset(r.GetPage(), r.GetPageSize())
}

type ExportRequest struct {
	QueryOptions
	Format  *string  `query:"format" json:"format"`
	Columns []string `query:"columns" json:"columns"`
}

func (r ExportRequest) GetFormat() string {
	value := strings.ToLower(helper.StringValue(r.Format))

	switch value {
	case ExportFormatCSV:
		return ExportFormatCSV
	case ExportFormatExcel:
		return ExportFormatExcel
	default:
		return DefaultExportFormat
	}
}
