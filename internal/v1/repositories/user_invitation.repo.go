package repositories

import (
	"github.com/MarcelArt/kas-bon-v2/internal/v1/models"
	"github.com/gofiber/fiber/v3"
	"github.com/morkid/paginate"
	"gorm.io/gorm"
)

type IUserInvitationRepo interface {
	Create(invitation models.UserInvitationInput) (uint, error)
	Read(c fiber.Ctx) (paginate.Page, []models.UserInvitationPage)
	ReadByUserID(c fiber.Ctx, userID any) (paginate.Page, []models.UserInvitationPage)
	Update(id any, invitation models.UserInvitationInput) error
	Delete(id any) error
	GetByID(id any) (models.UserInvitation, error)
}

type UserInvitationRepo struct {
	db        *gorm.DB
	pageQuery string
}

func NewUserInvitationRepo(db *gorm.DB) *UserInvitationRepo {
	return &UserInvitationRepo{
		db: db,
		pageQuery: `
			select 
				uv.*,
				u.username user,
				d.name domain,
				r.name role
			from user_invitations uv
			join users u on uv.user_id = u.id
			join domains d on uv.domain_id = d.id
			join roles r on uv.role_id = r.id
			where uv.deleted_at isnull
		`,
	}
}

func (r *UserInvitationRepo) Create(invitation models.UserInvitationInput) (uint, error) {
	err := r.db.Create(&invitation).Error
	return invitation.ID, err
}

func (r *UserInvitationRepo) Read(c fiber.Ctx) (paginate.Page, []models.UserInvitationPage) {
	var invitations []models.UserInvitationPage
	pg := paginate.New()

	stmt := r.db.Raw(r.pageQuery)

	page := pg.With(stmt).Request(c.Request()).Response(&invitations)
	return page, invitations
}

func (r *UserInvitationRepo) ReadByUserID(c fiber.Ctx, userID any) (paginate.Page, []models.UserInvitationPage) {
	var invitations []models.UserInvitationPage
	pg := paginate.New()

	stmt := r.db.Raw(r.pageQuery+" and uv.user_id = ?", userID)

	page := pg.With(stmt).Request(c.Request()).Response(&invitations)
	return page, invitations
}

func (r *UserInvitationRepo) Update(id any, invitation models.UserInvitationInput) error {
	return r.db.Model(&models.UserInvitation{}).Where("id = ?", id).Updates(&invitation).Error
}

func (r *UserInvitationRepo) Delete(id any) error {
	return r.db.Delete(&models.UserInvitation{}, id).Error
}

func (r *UserInvitationRepo) GetByID(id any) (models.UserInvitation, error) {
	var invitation models.UserInvitation
	err := r.db.Where("id = ?", id).
		Preload("User").
		Preload("Domain").
		Preload("Role").
		First(&invitation).Error
	return invitation, err
}
