package model

type Pagination struct {
	Page      int `json:"page"`
	PerPage   int `json:"per_page"`
	Total     int `json:"total"`
	TotalPages int `json:"total_pages"`
}

func NewPagination(page, perPage, total int) Pagination {
	totalPages := 0
	if perPage > 0 && total > 0 {
		totalPages = (total + perPage - 1) / perPage
	}
	return Pagination{
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	}
}

func (p Pagination) Offset() int {
	if p.Page <= 0 {
		return 0
	}
	return (p.Page - 1) * p.PerPage
}
