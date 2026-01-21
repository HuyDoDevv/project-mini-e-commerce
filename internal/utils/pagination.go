package utils

import (
	"strconv"
)

type Pagination struct {
	Page      int32 `json:"page"`
	Limit     int32 `json:"limit"`
	Total     int32 `json:"total"`
	TotalPage int32 `json:"total_page"`
	NextPage  bool  `json:"next_page"`
	PrevPage  bool  `json:"prev_page"`
}

func NewPagination(limit, page int32, totalRecords int64) *Pagination {
	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		envLimit := GetEnv("LIMIT_ITEM_ON_PER_PAGE", "10")
		limitInt, err := strconv.Atoi(envLimit)
		if err != nil || limitInt <= 0 {
			limit = 10
		} else {
			limit = int32(limitInt)
		}
	}

	totalPage := (int32(totalRecords) + limit - 1) / limit
	return &Pagination{
		Page:      page,
		Limit:     limit,
		Total:     int32(totalRecords),
		TotalPage: totalPage,
		NextPage:  page < totalPage,
		PrevPage:  page > 1,
	}
}

func NewPaginationResponse(data any, limit, page int32, totalRecords int64) map[string]any {
	return map[string]any{
		"data":       data,
		"pagination": NewPagination(limit, page, totalRecords),
	}
}
