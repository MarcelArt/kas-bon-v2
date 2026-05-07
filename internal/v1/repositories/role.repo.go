package repositories

import (
	"github.com/MarcelArt/kas-bon-v2/internal/v1/models"
	"github.com/gofiber/fiber/v3"
	"github.com/morkid/paginate"
	"gorm.io/gorm"
)

type IRoleRepo interface {
	Create(role models.RoleInput) (uint, error)
	Read(c fiber.Ctx, domainID any) (paginate.Page, []models.Role)
	Update(id any, role models.Role) error
	Delete(id any) error
	GetByID(id any) (models.Role, error)
}

type RoleRepo struct {
	db        *gorm.DB
	pageQuery string
}

func NewRoleRepo(db *gorm.DB) *RoleRepo {
	return &RoleRepo{
		db: db,
		pageQuery: `
			select * from roles where deleted_at isnull and domain_id = ?
		`,
	}
}

func (r *RoleRepo) Create(role models.RoleInput) (uint, error) {
	err := r.db.Create(&role).Error
	return role.ID, err
}

func (r *RoleRepo) Read(c fiber.Ctx, domainID any) (paginate.Page, []models.Role) {
	var roles []models.Role
	pg := paginate.New()

	stmt := r.db.Raw(r.pageQuery, domainID)

	page := pg.With(stmt).Request(c.Request()).Response(&roles)
	return page, roles
}

func (r *RoleRepo) Update(id any, role models.Role) error {
	return r.db.Model(&models.Role{}).Where("id = ?", id).Updates(&role).Error
}

func (r *RoleRepo) Delete(id any) error {
	return r.db.Delete(&models.Role{}, id).Error
}

func (r *RoleRepo) GetByID(id any) (models.Role, error) {
	var role models.Role
	err := r.db.Where("id = ?", id).First(&role).Error
	return role, err
}
