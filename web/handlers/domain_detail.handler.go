package handlers

import (
	"github.com/MarcelArt/kas-bon-v2/internal/v1/models"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/services"
	webModels "github.com/MarcelArt/kas-bon-v2/web/models"
	"github.com/gofiber/fiber/v3"
)

type DomainDetailHandler struct {
	domainSvc services.IDomainService
	roleSvc   services.IRoleService
}

func NewDomainDetailHandler(domainSvc services.IDomainService, roleSvc services.IRoleService) *DomainDetailHandler {
	return &DomainDetailHandler{domainSvc: domainSvc, roleSvc: roleSvc}
}

func (h *DomainDetailHandler) DomainDetailPage(c fiber.Ctx) error {
	domainID := c.Params("id")

	domain, err := h.domainSvc.GetByID(domainID)
	if err != nil {
		return c.Redirect().To("/domains")
	}

	page, roles := h.roleSvc.Read(c, domainID)

	viewDomain := webModels.DomainViewModel{
		ID:             domain.ID,
		Name:           domain.Name,
		Description:    domain.Description,
		IsOrganization: domain.IsOrganization,
		CreatedAt:      domain.CreatedAt,
	}

	viewRoles := make([]webModels.RoleViewModel, len(roles))
	perms := getPermissions(c)
	for i, r := range roles {
		viewRoles[i] = webModels.RoleViewModel{
			ID:                 r.ID,
			Name:               r.Name,
			Description:        r.Description,
			CreatedAt:          r.CreatedAt,
			CanUpdate:          perms["roles#update"],
			CanDelete:          perms["roles#delete"],
			CanReadPermissions: perms["roles#read"],
		}
	}

	domainUsers, _ := h.domainSvc.GetUsers(domainID)
	viewUsers := make([]webModels.DomainUserViewModel, len(domainUsers))
	for i, du := range domainUsers {
		roleName := ""
		if len(du.Policy) > 1 {
			roleName = du.Policy[1]
		}
		viewUsers[i] = webModels.DomainUserViewModel{
			Username:  du.User.Username,
			Email:     du.User.Email,
			RoleName:  roleName,
			CreatedAt: du.User.CreatedAt,
		}
	}

	basePath := "/domains/" + domainID
	data := webModels.DomainDetailPageData{
		PageData:   newPageData(c, domain.Name, "domain_detail"),
		Domain: viewDomain,
		Roles:  viewRoles,
		Users:  viewUsers,
		Pagination: webModels.NewPaginationData(
			page.Page, page.Size, page.TotalPages, page.Total,
			page.First, page.Last, basePath,
		),
	}

	if isHtmx(c) {
		return c.Render("domain_roles_table", data)
	}

	return c.Render("domain_detail", data)
}

func (h *DomainDetailHandler) CreateRoleForm(c fiber.Ctx) error {
	domainID := c.Params("id")
	return c.Render("role_form_create", map[string]string{"DomainID": domainID})
}

func (h *DomainDetailHandler) EditRoleForm(c fiber.Ctx) error {
	roleID := c.Params("id")

	role, err := h.roleSvc.GetByID(roleID)
	if err != nil {
		return c.Redirect().To("/domains")
	}

	viewRole := webModels.RoleViewModel{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		CreatedAt:   role.CreatedAt,
	}

	return c.Render("role_form_edit", viewRole)
}

func (h *DomainDetailHandler) CreateRole(c fiber.Ctx) error {
	domainID := c.Params("id")

	var input models.RoleInput
	if err := c.Bind().Form(&input); err != nil {
		return c.Redirect().To("/domains/" + domainID)
	}
	input.DomainID = parseUint(domainID)

	id, err := h.roleSvc.Create(input)
	if err != nil {
		return c.Redirect().To("/domains/" + domainID)
	}

	if isHtmx(c) {
		return h.renderRoleRow(c, id)
	}

	return c.Redirect().To("/domains/" + domainID)
}

func (h *DomainDetailHandler) UpdateRole(c fiber.Ctx) error {
	roleID := c.Params("id")

	var input models.Role
	if err := c.Bind().Form(&input); err != nil {
		return c.Redirect().To("/domains")
	}

	if err := h.roleSvc.Update(roleID, input); err != nil {
		return c.Redirect().To("/domains")
	}

	if isHtmx(c) {
		return h.renderRoleRow(c, roleID)
	}

	role, _ := h.roleSvc.GetByID(roleID)
	return c.Redirect().To("/domains/" + uintToString(role.DomainID))
}

func (h *DomainDetailHandler) DeleteRole(c fiber.Ctx) error {
	roleID := c.Params("id")

	if err := h.roleSvc.Delete(roleID); err != nil {
		return c.Redirect().To("/domains")
	}

	if isHtmx(c) {
		return c.SendStatus(fiber.StatusOK)
	}

	return c.Redirect().To("/domains")
}

func (h *DomainDetailHandler) renderRoleRow(c fiber.Ctx, id any) error {
	role, err := h.roleSvc.GetByID(id)
	if err != nil {
		return c.Redirect().To("/domains")
	}

	perms := getPermissions(c)
	viewRole := webModels.RoleViewModel{
		ID:                 role.ID,
		Name:               role.Name,
		Description:        role.Description,
		CreatedAt:          role.CreatedAt,
		CanUpdate:          perms["roles#update"],
		CanDelete:          perms["roles#delete"],
		CanReadPermissions: perms["roles#read"],
	}

	return c.Render("role_row", viewRole)
}

func parseUint(s string) uint {
	var n uint
	for _, c := range s {
		if c < '0' || c > '9' {
			break
		}
		n = n*10 + uint(c-'0')
	}
	return n
}

func uintToString(n uint) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte(n%10) + '0'
		n /= 10
	}
	return string(buf[i:])
}
