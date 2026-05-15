package handlers

import (
	"strconv"

	"github.com/MarcelArt/kas-bon-v2/internal/v1/models"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/services"
	webModels "github.com/MarcelArt/kas-bon-v2/web/models"
	"github.com/gofiber/fiber/v3"
)

type DomainHandler struct {
	domainSvc services.IDomainService
}

func NewDomainHandler(domainSvc services.IDomainService) *DomainHandler {
	return &DomainHandler{domainSvc: domainSvc}
}

func (h *DomainHandler) DomainsPage(c fiber.Ctx) error {
	userID := c.Locals("userID")
	parentID, _ := strconv.ParseUint(c.Cookies("current_domain_id"), 10, 64)
	page, domains := h.domainSvc.GetUserDomains(c, userID, uint(parentID))
	perms := getPermissions(c)

	viewDomains := make([]webModels.DomainViewModel, len(domains))
	for i, d := range domains {
		viewDomains[i] = webModels.DomainViewModel{
			ID:             d.ID,
			Name:           d.Name,
			Description:    d.Description,
			IsOrganization: d.IsOrganization,
			CreatedAt:      d.CreatedAt,
			CanUpdate:      perms["domains#update"],
			CanDelete:      perms["domains#delete"],
		}
	}

	data := webModels.DomainsPageData{
		PageData:   newPageData(c, "Domains", "domains"),
		Domains: viewDomains,
		Pagination: webModels.NewPaginationData(
			page.Page, page.Size, page.TotalPages, page.Total,
			page.First, page.Last, "/domains",
		),
	}

	if isHtmx(c) {
		return c.Render("domains_table", data)
	}

	return c.Render("domains", data)
}

func (h *DomainHandler) CreateDomainForm(c fiber.Ctx) error {
	return c.Render("domain_form_create", nil)
}

func (h *DomainHandler) EditDomainForm(c fiber.Ctx) error {
	id := c.Params("id")

	domain, err := h.domainSvc.GetByID(id)
	if err != nil {
		return c.Redirect().To("/domains")
	}

	viewDomain := webModels.DomainViewModel{
		ID:             domain.ID,
		Name:           domain.Name,
		Description:    domain.Description,
		IsOrganization: domain.IsOrganization,
		CreatedAt:      domain.CreatedAt,
	}

	return c.Render("domain_form_edit", viewDomain)
}

func (h *DomainHandler) CreateDomain(c fiber.Ctx) error {
	var input models.DomainInput
	if err := c.Bind().Form(&input); err != nil {
		return c.Redirect().To("/domains")
	}

	id, err := h.domainSvc.Create(input, c.Locals("userID"))
	if err != nil {
		return c.Redirect().To("/domains")
	}

	if isHtmx(c) {
		return h.renderDomainRow(c, id)
	}

	return c.Redirect().To("/domains")
}

func (h *DomainHandler) UpdateDomain(c fiber.Ctx) error {
	id := c.Params("id")

	var input models.Domain
	if err := c.Bind().Form(&input); err != nil {
		return c.Redirect().To("/domains")
	}

	if err := h.domainSvc.Update(id, input); err != nil {
		return c.Redirect().To("/domains")
	}

	if isHtmx(c) {
		return h.renderDomainRow(c, id)
	}

	return c.Redirect().To("/domains")
}

func (h *DomainHandler) DeleteDomain(c fiber.Ctx) error {
	id := c.Params("id")

	if err := h.domainSvc.Delete(id); err != nil {
		return c.Redirect().To("/domains")
	}

	if isHtmx(c) {
		return c.SendStatus(fiber.StatusOK)
	}

	return c.Redirect().To("/domains")
}

func (h *DomainHandler) renderDomainRow(c fiber.Ctx, id any) error {
	domain, err := h.domainSvc.GetByID(id)
	if err != nil {
		return c.Redirect().To("/domains")
	}

	perms := getPermissions(c)
	viewDomain := webModels.DomainViewModel{
		ID:             domain.ID,
		Name:           domain.Name,
		Description:    domain.Description,
		IsOrganization: domain.IsOrganization,
		CreatedAt:      domain.CreatedAt,
		CanUpdate:      perms["domains#update"],
		CanDelete:      perms["domains#delete"],
	}

	return c.Render("domain_row", viewDomain)
}
