package handlers

import (
	"fmt"

	"github.com/MarcelArt/kas-bon-v2/internal/v1/models"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/services"
	webModels "github.com/MarcelArt/kas-bon-v2/web/models"
	"github.com/gofiber/fiber/v3"
)

type AppHandler struct {
	appSvc services.IAppService
}

func NewAppHandler(appSvc services.IAppService) *AppHandler {
	return &AppHandler{appSvc: appSvc}
}

func (h *AppHandler) AppsPage(c fiber.Ctx) error {
	page, apps := h.appSvc.Read(c)

	viewApps := make([]webModels.AppViewModel, len(apps))
	for i, a := range apps {
		viewApps[i] = webModels.AppViewModel{
			ID:          a.ID,
			Name:        a.Name,
			Description: a.Description,
			CreatedAt:   a.CreatedAt,
		}
	}

	data := webModels.AppsPageData{
		PageData: webModels.PageData{
			Title:      "Apps",
			ActivePage: "apps",
		},
		Apps: viewApps,
		Pagination: webModels.NewPaginationData(
			page.Page, page.Size, page.TotalPages, page.Total,
			page.First, page.Last, "/apps",
		),
	}

	if isHtmx(c) {
		return c.Render("apps_table", data)
	}

	return c.Render("apps", data)
}

func (h *AppHandler) CreateAppForm(c fiber.Ctx) error {
	return c.Render("app_form_create", nil)
}

func (h *AppHandler) EditAppForm(c fiber.Ctx) error {
	id := c.Params("id")

	app, err := h.appSvc.GetByID(id)
	if err != nil {
		return c.Redirect().To("/apps")
	}

	viewApp := webModels.AppViewModel{
		ID:          app.ID,
		Name:        app.Name,
		Description: app.Description,
		CreatedAt:   app.CreatedAt,
	}

	return c.Render("app_form_edit", viewApp)
}

func (h *AppHandler) CreateApp(c fiber.Ctx) error {
	var input models.AppInput
	if err := c.Bind().Form(&input); err != nil {
		return c.Redirect().To("/apps")
	}

	id, err := h.appSvc.Create(input)
	if err != nil {
		return c.Redirect().To("/apps")
	}

	if isHtmx(c) {
		return h.renderAppRow(c, id)
	}

	return c.Redirect().To("/apps")
}

func (h *AppHandler) UpdateApp(c fiber.Ctx) error {
	id := c.Params("id")

	var input models.App
	if err := c.Bind().Form(&input); err != nil {
		return c.Redirect().To("/apps")
	}

	if err := h.appSvc.Update(id, input); err != nil {
		return c.Redirect().To("/apps")
	}

	if isHtmx(c) {
		return h.renderAppRow(c, id)
	}

	return c.Redirect().To("/apps")
}

func (h *AppHandler) DeleteApp(c fiber.Ctx) error {
	id := c.Params("id")

	if err := h.appSvc.Delete(id); err != nil {
		return c.Redirect().To("/apps")
	}

	if isHtmx(c) {
		return c.SendStatus(fiber.StatusOK)
	}

	return c.Redirect().To("/apps")
}

func (h *AppHandler) renderAppRow(c fiber.Ctx, id any) error {
	app, err := h.appSvc.GetByID(id)
	if err != nil {
		return c.Redirect().To("/apps")
	}

	viewApp := webModels.AppViewModel{
		ID:          app.ID,
		Name:        app.Name,
		Description: app.Description,
		CreatedAt:   app.CreatedAt,
	}

	return c.Render("app_row", viewApp)
}

func isHtmx(c fiber.Ctx) bool {
	return c.Get("HX-Request") == "true"
}

func renderWebAlert(c fiber.Ctx, alertType, message string) error {
	html := fmt.Sprintf(`<div class="alert alert-%s">%s</div>`, alertType, message)
	return c.SendString(html)
}
