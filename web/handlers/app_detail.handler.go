package handlers

import (
	"github.com/MarcelArt/kas-bon-v2/internal/v1/models"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/services"
	webModels "github.com/MarcelArt/kas-bon-v2/web/models"
	"github.com/gofiber/fiber/v3"
)

type AppDetailHandler struct {
	appSvc  services.IAppService
	permSvc services.IPermissionService
}

func NewAppDetailHandler(appSvc services.IAppService, permSvc services.IPermissionService) *AppDetailHandler {
	return &AppDetailHandler{appSvc: appSvc, permSvc: permSvc}
}

func (h *AppDetailHandler) AppDetailPage(c fiber.Ctx) error {
	appID := c.Params("id")

	app, err := h.appSvc.GetByID(appID)
	if err != nil {
		return c.Redirect().To("/apps")
	}

	page, permissions := h.permSvc.Read(c, appID)
	perms := getPermissions(c)

	viewApp := webModels.AppViewModel{
		ID:          app.ID,
		Name:        app.Name,
		Description: app.Description,
		CreatedAt:   app.CreatedAt,
	}

	viewPerms := make([]webModels.PermissionViewModel, len(permissions))
	for i, p := range permissions {
		viewPerms[i] = webModels.PermissionViewModel{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			CreatedAt:   p.CreatedAt,
			CanUpdate:   perms["permissions#update"],
			CanDelete:   perms["permissions#delete"],
		}
	}

	basePath := "/apps/" + appID
	data := webModels.AppDetailPageData{
		PageData:       newPageData(c, app.Name, "app_detail"),
		App:            viewApp,
		PermissionList: viewPerms,
		Pagination: webModels.NewPaginationData(
			page.Page, page.Size, page.TotalPages, page.Total,
			page.First, page.Last, basePath,
		),
	}

	if isHtmx(c) {
		return c.Render("app_permissions_table", data)
	}

	return c.Render("app_detail", data)
}

func (h *AppDetailHandler) CreatePermissionForm(c fiber.Ctx) error {
	appID := c.Params("id")
	return c.Render("permission_form_create", map[string]string{"AppID": appID})
}

func (h *AppDetailHandler) EditPermissionForm(c fiber.Ctx) error {
	permID := c.Params("id")

	perm, err := h.permSvc.GetByID(permID)
	if err != nil {
		return c.Redirect().To("/apps")
	}

	viewPerm := webModels.PermissionViewModel{
		ID:          perm.ID,
		Name:        perm.Name,
		Description: perm.Description,
		CreatedAt:   perm.CreatedAt,
	}

	return c.Render("permission_form_edit", viewPerm)
}

func (h *AppDetailHandler) CreatePermission(c fiber.Ctx) error {
	appID := c.Params("id")

	var input models.PermissionInput
	if err := c.Bind().Form(&input); err != nil {
		return c.Redirect().To("/apps/" + appID)
	}
	input.AppID = parseUint(appID)

	id, err := h.permSvc.Create(input)
	if err != nil {
		return c.Redirect().To("/apps/" + appID)
	}

	if isHtmx(c) {
		return h.renderPermissionRow(c, id)
	}

	return c.Redirect().To("/apps/" + appID)
}

func (h *AppDetailHandler) UpdatePermission(c fiber.Ctx) error {
	permID := c.Params("id")

	var input models.Permission
	if err := c.Bind().Form(&input); err != nil {
		return c.Redirect().To("/apps")
	}

	if err := h.permSvc.Update(permID, input); err != nil {
		return c.Redirect().To("/apps")
	}

	if isHtmx(c) {
		return h.renderPermissionRow(c, permID)
	}

	perm, _ := h.permSvc.GetByID(permID)
	return c.Redirect().To("/apps/" + uintToString(perm.AppID))
}

func (h *AppDetailHandler) DeletePermission(c fiber.Ctx) error {
	permID := c.Params("id")

	if err := h.permSvc.Delete(permID); err != nil {
		return c.Redirect().To("/apps")
	}

	if isHtmx(c) {
		return c.SendStatus(fiber.StatusOK)
	}

	return c.Redirect().To("/apps")
}

func (h *AppDetailHandler) renderPermissionRow(c fiber.Ctx, id any) error {
	perm, err := h.permSvc.GetByID(id)
	if err != nil {
		return c.Redirect().To("/apps")
	}

	p := getPermissions(c)
	viewPerm := webModels.PermissionViewModel{
		ID:          perm.ID,
		Name:        perm.Name,
		Description: perm.Description,
		CreatedAt:   perm.CreatedAt,
		CanUpdate:   p["permissions#update"],
		CanDelete:   p["permissions#delete"],
	}

	return c.Render("permission_row", viewPerm)
}
