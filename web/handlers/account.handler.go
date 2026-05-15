package handlers

import (
	"strconv"

	"github.com/MarcelArt/kas-bon-v2/internal/v1/models"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/services"
	webModels "github.com/MarcelArt/kas-bon-v2/web/models"
	"github.com/alexedwards/argon2id"
	"github.com/gofiber/fiber/v3"
)

type AccountHandler struct {
	userSvc       services.IUserService
	domainSvc     services.IDomainService
	invitationSvc services.IUserInvitationService
}

func NewAccountHandler(userSvc services.IUserService, domainSvc services.IDomainService, invitationSvc services.IUserInvitationService) *AccountHandler {
	return &AccountHandler{userSvc: userSvc, domainSvc: domainSvc, invitationSvc: invitationSvc}
}

func (h *AccountHandler) AccountPage(c fiber.Ctx) error {
	userID := c.Locals("userID")

	user, err := h.userSvc.GetByID(userID)
	if err != nil {
		return c.Redirect().To("/login")
	}

	if domID, _ := strconv.ParseUint(c.Cookies("current_domain_id"), 10, 64); domID > 0 {
		dom, err := h.domainSvc.GetByID(domID)
		if err == nil {
			c.Locals("currentOrgName", dom.Name)
			c.Locals("currentOrgID", uint(domID))
		}
	}

	page, invitations := h.invitationSvc.ReadByUserID(c, userID)

	viewInvitations := make([]webModels.InvitationViewModel, len(invitations))
	for i, inv := range invitations {
		viewInvitations[i] = webModels.InvitationViewModel{
			ID:         inv.ID,
			Domain:     inv.Domain,
			Role:       inv.Role,
			CreatedAt:  inv.CreatedAt,
			AcceptedAt: inv.AcceptedAt,
			RejectedAt: inv.RejectedAt,
		}
	}

	data := webModels.AccountPageData{
		PageData: newPageData(c, "Account", "account"),
		User: webModels.UserViewModel{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		},
		Invitations: viewInvitations,
		Pagination: webModels.NewPaginationData(
			page.Page, page.Size, page.TotalPages, page.Total,
			page.First, page.Last, "/account",
		),
	}

	if isHtmx(c) {
		return c.Render("account_invitations_table", data)
	}

	return c.Render("account", data)
}

func (h *AccountHandler) UpdateAccount(c fiber.Ctx) error {
	userID := c.Locals("userID")

	var form webModels.UpdateAccountForm
	if err := c.Bind().Form(&form); err != nil {
		return renderAlert(c, "error", "Invalid form data")
	}

	updates := models.User{
		Email: form.Email,
	}

	if form.Password != "" {
		hash, err := argon2id.CreateHash(form.Password, argon2id.DefaultParams)
		if err != nil {
			return renderAlert(c, "error", "Failed to update password")
		}
		updates.Password = hash
	}

	if err := h.userSvc.Update(userID, updates); err != nil {
		return renderAlert(c, "error", "Failed to update account")
	}

	return renderAlert(c, "success", "Account updated successfully")
}
