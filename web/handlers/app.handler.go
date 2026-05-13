package handlers

import (
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

	return c.Render("apps", webModels.AppsPageData{
		PageData: webModels.PageData{
			Title:      "Apps",
			ActivePage: "apps",
		},
		Apps: viewApps,
	})
}
