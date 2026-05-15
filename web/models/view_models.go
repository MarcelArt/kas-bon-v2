package models

import (
	"strings"
	"time"
)

type PageData struct {
	Title       string
	ActivePage  string
	Permissions map[string]bool
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
	LastPage   int64
	BasePath   string
	TargetID   string
}

func NewPaginationData(page, size, totalPages, total int64, first, last bool, basePath string) PaginationData {
	prev := page - 1
	if prev < 0 {
		prev = 0
	}
	maxPage := totalPages - 1
	if maxPage < 0 {
		maxPage = 0
	}
	next := page + 1
	if next > maxPage {
		next = maxPage
	}
	return PaginationData{
		Page:       page + 1,
		Size:       size,
		TotalPages: totalPages,
		Total:      total,
		First:      first,
		Last:       last,
		PrevPage:   prev,
		NextPage:   next,
		LastPage:   maxPage,
		BasePath:   basePath,
		TargetID:   strings.ReplaceAll(basePath[1:], "/", "-") + "-table-area",
	}
}

type AppViewModel struct {
	ID          uint
	Name        string
	Description string
	CreatedAt   time.Time
	CanUpdate   bool
	CanDelete   bool
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
	CanUpdate      bool
	CanDelete      bool
}

type DomainsPageData struct {
	PageData
	Domains    []DomainViewModel
	Pagination PaginationData
}

type DomainUserViewModel struct {
	ID        uint
	Username  string
	Email     string
	RoleNames []string
	CreatedAt time.Time
}

type DomainDetailPageData struct {
	PageData
	Domain     DomainViewModel
	Roles      []RoleViewModel
	Users      []DomainUserViewModel
	Pagination PaginationData
}

type RoleViewModel struct {
	ID                  uint
	Name                string
	Description         string
	CreatedAt           time.Time
	IsAssigned          bool
	CanUpdate           bool
	CanDelete           bool
	CanReadPermissions  bool
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

type PermissionViewModel struct {
	ID          uint
	Name        string
	Description string
	CreatedAt   time.Time
	IsAssigned  bool
	CanUpdate   bool
	CanDelete   bool
}

type AppDetailPageData struct {
	PageData
	App            AppViewModel
	PermissionList []PermissionViewModel
	Pagination     PaginationData
}

type RolePermissionsPageData struct {
	PageData
	Role     RoleViewModel
	DomainID uint
	Apps     []AppViewModel
}

type RolePermissionsListData struct {
	RoleID      uint
	AppID       uint
	Permissions []PermissionViewModel
}

type AlertData struct {
	Type    string
	Message string
}

type OrgViewModel struct {
	ID          uint
	Name        string
	Description string
}

type OrgSelectPageData struct {
	PageData
	Organizations []OrgViewModel
}

type InviteUserFormData struct {
	DomainID uint
	Users    []UserOption
	Roles    []RoleOption
}

type UserOption struct {
	ID       uint
	Username string
	Email    string
}

type RoleOption struct {
	ID   uint
	Name string
}

type UserRolesPageData struct {
	PageData
	User       UserViewModel
	DomainID   uint
	DomainName string
	Roles      []RoleViewModel
}

type UserViewModel struct {
	ID       uint
	Username string
	Email    string
}
