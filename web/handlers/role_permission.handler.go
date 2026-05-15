package handlers

import (
	"github.com/MarcelArt/kas-bon-v2/internal/v1/services"
	webModels "github.com/MarcelArt/kas-bon-v2/web/models"
	"github.com/gofiber/fiber/v3"
)

type RolePermissionHandler struct {
	roleSvc services.IRoleService
	permSvc services.IPermissionService
	appSvc  services.IAppService
}

func NewRolePermissionHandler(roleSvc services.IRoleService, permSvc services.IPermissionService, appSvc services.IAppService) *RolePermissionHandler {
	return &RolePermissionHandler{roleSvc: roleSvc, permSvc: permSvc, appSvc: appSvc}
}

func (h *RolePermissionHandler) PermissionsPage(c fiber.Ctx) error {
	roleID := c.Params("id")

	role, err := h.roleSvc.GetByID(roleID)
	if err != nil {
		return c.Redirect().To("/domains")
	}

	_, apps := h.appSvc.Read(c)

	viewApps := make([]webModels.AppViewModel, len(apps))
	for i, a := range apps {
		viewApps[i] = webModels.AppViewModel{
			ID:          a.ID,
			Name:        a.Name,
			Description: a.Description,
			CreatedAt:   a.CreatedAt,
		}
	}

	data := webModels.RolePermissionsPageData{
		PageData:   newPageData(c, "Domains", "role_permissions"),
		Role: webModels.RoleViewModel{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			CreatedAt:   role.CreatedAt,
		},
		DomainID: role.DomainID,
		Apps:     viewApps,
	}

	return c.Render("role_permissions", data)
}

func (h *RolePermissionHandler) PermissionsList(c fiber.Ctx) error {
	roleID := c.Params("id")
	appID := c.Query("app_id")

	if appID == "" {
		return c.SendString("")
	}

	return h.renderPermissionsList(c, roleID, parseUint(appID))
}

func (h *RolePermissionHandler) AssignPermissions(c fiber.Ctx) error {
	roleID := c.Params("id")
	appID := parseUint(c.FormValue("app_id"))

	if appID == 0 {
		return c.SendString("")
	}

	var permissionIDs []uint
	c.Request().PostArgs().VisitAll(func(key, value []byte) {
		if string(key) == "permission_ids" {
			if id := parseUint(string(value)); id > 0 {
				permissionIDs = append(permissionIDs, id)
			}
		}
	})

	h.roleSvc.AssignPermissions(roleID, appID, permissionIDs)

	return h.renderPermissionsList(c, roleID, appID)
}

func (h *RolePermissionHandler) renderPermissionsList(c fiber.Ctx, roleID string, appID uint) error {
	permissions, err := h.permSvc.GetByAppID(appID)
	if err != nil {
		return c.SendString("")
	}

	assignedSet := make(map[uint]bool)
	assignedResp, _ := h.roleSvc.GetPermissions(roleID)
	if assignedResp.Permissions != nil {
		for _, p := range assignedResp.Permissions {
			if p.AppID == appID {
				assignedSet[p.ID] = true
			}
		}
	}

	viewPerms := make([]webModels.PermissionViewModel, len(permissions))
	for i, p := range permissions {
		viewPerms[i] = webModels.PermissionViewModel{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			IsAssigned:  assignedSet[p.ID],
		}
	}

	data := webModels.RolePermissionsListData{
		RoleID:      parseUint(roleID),
		AppID:       appID,
		Permissions: viewPerms,
	}

	return c.Render("role_permissions_list", data)
}
