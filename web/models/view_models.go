package models

import "time"

type PageData struct {
	Title      string
	ActivePage string
}

type PaginationData struct {
	Page       int64
	Size       int64
	TotalPages int64
	Total      int64
	First      bool
	Last       bool
	PrevPage   int64
	NextPage   int64
	BasePath   string
	TargetID   string
}

func NewPaginationData(page, size, totalPages, total int64, first, last bool, basePath string) PaginationData {
	prev := page - 1
	if prev < 1 {
		prev = 1
	}
	next := page + 1
	if next > totalPages {
		next = totalPages
	}
	return PaginationData{
		Page:       page,
		Size:       size,
		TotalPages: totalPages,
		Total:      total,
		First:      first,
		Last:       last,
		PrevPage:   prev,
		NextPage:   next,
		BasePath:   basePath,
		TargetID:   basePath[1:] + "-table-area",
	}
}

type AppViewModel struct {
	ID          uint
	Name        string
	Description string
	CreatedAt   time.Time
}

type AppsPageData struct {
	PageData
	Apps       []AppViewModel
	Pagination PaginationData
}

type DomainViewModel struct {
	ID             uint
	Name           string
	Description    string
	IsOrganization bool
	CreatedAt      time.Time
}

type DomainsPageData struct {
	PageData
	Domains    []DomainViewModel
	Pagination PaginationData
}

type DomainDetailPageData struct {
	PageData
	Domain     DomainViewModel
	Roles      []RoleViewModel
	Pagination PaginationData
}

type RoleViewModel struct {
	ID          uint
	Name        string
	Description string
	CreatedAt   time.Time
}

type LoginForm struct {
	Username   string `form:"username"`
	Password   string `form:"password"`
	IsRemember bool   `form:"isRemember"`
}

type RegisterForm struct {
	Username        string `form:"username"`
	Email           string `form:"email"`
	Password        string `form:"password"`
	ConfirmPassword string `form:"confirmPassword"`
}

type AlertData struct {
	Type    string
	Message string
}
