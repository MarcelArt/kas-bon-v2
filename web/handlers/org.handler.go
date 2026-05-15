package handlers

import (
	"strconv"

	"github.com/MarcelArt/kas-bon-v2/internal/v1/models"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/services"
	webModels "github.com/MarcelArt/kas-bon-v2/web/models"
	"github.com/gofiber/fiber/v3"
)

type OrgHandler struct {
	domainSvc services.IDomainService
	userSvc   services.IUserService
}

func NewOrgHandler(domainSvc services.IDomainService, userSvc services.IUserService) *OrgHandler {
	return &OrgHandler{domainSvc: domainSvc, userSvc: userSvc}
}

func (h *OrgHandler) OrgsPage(c fiber.Ctx) error {
	userID := c.Locals("userID")
	currentOrgID, _ := strconv.ParseUint(c.Cookies("current_domain_id"), 10, 64)

	orgs, err := h.userSvc.GetOrganizations(userID)
	if err != nil {
		return c.Redirect().To("/apps")
	}

	viewOrgs := make([]webModels.OrgViewModel, len(orgs))
	for i, o := range orgs {
		viewOrgs[i] = webModels.OrgViewModel{
			ID:           o.ID,
			Name:         o.Name,
			Description:  o.Description,
			IsCurrentOrg: o.ID == uint(currentOrgID),
		}
	}

	perms := getPermissions(c)
	data := webModels.OrganizationsPageData{
		PageData:      newPageData(c, "Organizations", "organizations"),
		Organizations: viewOrgs,
	}

	if isHtmx(c) {
		data.Permissions = perms
		return c.Render("orgs_table", data)
	}

	return c.Render("organizations", data)
}

func (h *OrgHandler) SwitchOrg(c fiber.Ctx) error {
	domainID := c.FormValue("domain_id")
	if domainID == "" {
		return c.Redirect().To("/organizations")
	}

	id := parseUint(domainID)
	if id == 0 {
		return c.Redirect().To("/organizations")
	}

	setDomainCookie(c, id)

	c.Set("HX-Redirect", "/domains")
	return c.SendStatus(fiber.StatusOK)
}

func (h *OrgHandler) CreateOrgForm(c fiber.Ctx) error {
	return c.Render("org_form_create", nil)
}

func (h *OrgHandler) CreateOrg(c fiber.Ctx) error {
	var input models.DomainInput
	if err := c.Bind().Form(&input); err != nil {
		return c.Redirect().To("/organizations")
	}
	input.IsOrganization = true

	id, err := h.domainSvc.Create(input, c.Locals("userID"))
	if err != nil {
		return c.Redirect().To("/organizations")
	}

	if isHtmx(c) {
		return h.renderOrgRow(c, id)
	}

	return c.Redirect().To("/organizations")
}

func (h *OrgHandler) renderOrgRow(c fiber.Ctx, id any) error {
	domain, err := h.domainSvc.GetByID(id)
	if err != nil {
		return c.Redirect().To("/organizations")
	}

	currentOrgID, _ := strconv.ParseUint(c.Cookies("current_domain_id"), 10, 64)

	viewOrg := webModels.OrgViewModel{
		ID:           domain.ID,
		Name:         domain.Name,
		Description:  domain.Description,
		IsCurrentOrg: domain.ID == uint(currentOrgID),
	}

	return c.Render("org_row", viewOrg)
}
