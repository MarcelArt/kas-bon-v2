package handlers

import (
	"github.com/MarcelArt/kas-bon-v2/internal/v1/services"
	webModels "github.com/MarcelArt/kas-bon-v2/web/models"
	"github.com/gofiber/fiber/v3"
)

type UserRoleHandler struct {
	userSvc services.IUserService
	roleSvc services.IRoleService
	domSvc  services.IDomainService
}

func NewUserRoleHandler(userSvc services.IUserService, roleSvc services.IRoleService, domSvc services.IDomainService) *UserRoleHandler {
	return &UserRoleHandler{userSvc: userSvc, roleSvc: roleSvc, domSvc: domSvc}
}

func (h *UserRoleHandler) UserRolesPage(c fiber.Ctx) error {
	userID := c.Params("userId")
	domainID := c.Params("domainId")

	user, err := h.userSvc.GetByID(userID)
	if err != nil {
		return c.Redirect().To("/domains")
	}

	domain, err := h.domSvc.GetByID(domainID)
	if err != nil {
		return c.Redirect().To("/domains")
	}

	roles, err := h.roleSvc.GetByDomainID(domainID)
	if err != nil {
		return c.Redirect().To("/domains/" + domainID)
	}

	assignedRoles, _ := h.userSvc.GetRoles(userID, domainID)
	assignedSet := make(map[string]bool)
	for _, r := range assignedRoles {
		assignedSet[r] = true
	}

	viewRoles := make([]webModels.RoleViewModel, len(roles))
	for i, r := range roles {
		viewRoles[i] = webModels.RoleViewModel{
			ID:          r.ID,
			Name:        r.Name,
			Description: r.Description,
			CreatedAt:   r.CreatedAt,
			IsAssigned:  assignedSet[r.Name],
		}
	}

	data := webModels.UserRolesPageData{
		PageData: newPageData(c, "Domains", "user_roles"),
		User: webModels.UserViewModel{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		},
		DomainID:   domain.ID,
		DomainName: domain.Name,
		Roles:      viewRoles,
	}

	if isHtmx(c) {
		return c.Render("user_roles_list", data)
	}

	return c.Render("user_roles", data)
}

func (h *UserRoleHandler) AssignRoles(c fiber.Ctx) error {
	userID := c.Params("userId")
	domainID := c.Params("domainId")

	var permissionIDs []uint
	c.Request().PostArgs().VisitAll(func(key, value []byte) {
		if string(key) == "role_ids" {
			if id := parseUint(string(value)); id > 0 {
				permissionIDs = append(permissionIDs, id)
			}
		}
	})

	h.userSvc.AssignRoles(userID, domainID, permissionIDs)

	return h.UserRolesPage(c)
}
