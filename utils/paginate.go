package utils

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"math"
	"strconv"
)

type Pagination struct {
	Limit int         `json:"limit,omitempty;query:limit"`
	Page  int         `json:"page,omitempty;query:page"`
	Sort  *string     `json:"sort,omitempty;query:sort"`
	Total int64       `json:"total"`
	Pages int         `json:"pages"`
	Rows  interface{} `json:"rows"`
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination) GetLimit() int {
	if p.Limit == 0 {
		p.Limit = 12
	}
	return p.Limit
}

func (p *Pagination) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *Pagination) GetSort() string {
	if p.Sort == nil {
		return ""
	}
	return *p.Sort
}

func Paginate(value interface{}, pagination *Pagination, query *gorm.DB) func(db *gorm.DB) *gorm.DB {
	var total int64
	query.Model(value).Count(&total)

	pagination.Total = total
	pages := int(math.Ceil(float64(total) / float64(pagination.GetLimit())))
	pagination.Pages = pages

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order(pagination.GetSort())
	}
}

func PaginationFromContext(ctx echo.Context) Pagination {
	pagination := new(Pagination)
	sort := ctx.QueryParam("sort")

	pagination.Limit, _ = strconv.Atoi(ctx.QueryParam("limit"))
	pagination.Page, _ = strconv.Atoi(ctx.QueryParam("page"))
	pagination.Sort = &sort

	return *pagination
}
