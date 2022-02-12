package grpc_tpl

const GlobalPagePageTpl = `package page

import "math"

const (
	DefaultPageSize = 10
	MaxPageSize     = 100
	DefaultPage     = 1
)

type Pagination struct {
	CurrentPage int   {{.tagName}}json:"current_page"{{.tagName}} // 当前页码
	PerPage     int   {{.tagName}}json:"per_page"{{.tagName}}     // 当前页行数
	TotalPage   int   {{.tagName}}json:"total_page"{{.tagName}}   // 总页码
	Total       int64 {{.tagName}}json:"total"{{.tagName}}        // 总行数
}

func GetPagination(page, pageSize int, total int64) Pagination {
	return Pagination{
		Total:       total,
		CurrentPage: page,
		PerPage:     pageSize,
		TotalPage:   int(math.Ceil(float64(total) / float64(pageSize))),
	}
}
`
