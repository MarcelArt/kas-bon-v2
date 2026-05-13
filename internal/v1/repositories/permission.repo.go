package repositories

import (
	"github.com/MarcelArt/kas-bon-v2/internal/v1/models"
	"github.com/gofiber/fiber/v3"
	"github.com/morkid/paginate"
	"gorm.io/gorm"
)

type IPermissionRepo interface {
	Create(permission models.PermissionInput) (uint, error)
	Read(c fiber.Ctx, appID any) (paginate.Page, []models.Permission)
	Update(id any, permission models.Permission) error
	Delete(id any) error
	GetByID(id any) (models.Permission, error)
	GetDistinctByIDs(ids []uint) ([]string, error)
	GetByAppID(appID any) ([]models.Permission, error)
	GetByNames(names []string) ([]models.Permission, error)
}

type PermissionRepo struct {
	db        *gorm.DB
	pageQuery string
}

func NewPermissionRepo(db *gorm.DB) *PermissionRepo {
	return &PermissionRepo{
		db: db,
		pageQuery: `
			select * from permissions where deleted_at isnull and app_id = ?
		`,
	}
}

func (r *PermissionRepo) Create(permission models.PermissionInput) (uint, error) {
	err := r.db.Create(&permission).Error
	return permission.ID, err
}

func (r *PermissionRepo) Read(c fiber.Ctx, appID any) (paginate.Page, []models.Permission) {
	var permissions []models.Permission
	pg := paginate.New()

	stmt := r.db.Raw(r.pageQuery, appID)

	page := pg.With(stmt).Request(c.Request()).Response(&permissions)
	return page, permissions
}

func (r *PermissionRepo) Update(id any, permission models.Permission) error {
	return r.db.Model(&models.Permission{}).Where("id = ?", id).Updates(&permission).Error
}

func (r *PermissionRepo) Delete(id any) error {
	return r.db.Delete(&models.Permission{}, id).Error
}

func (r *PermissionRepo) GetByID(id any) (models.Permission, error) {
	var permission models.Permission
	err := r.db.Where("id = ?", id).First(&permission).Error
	return permission, err
}

func (r *PermissionRepo) GetDistinctByIDs(ids []uint) ([]string, error) {
	var permissions []string
	err := r.db.Model(models.Permission{}).Where("id in ?", ids).Distinct("name").Pluck("name", &permissions).Error
	return permissions, err
}

func (r *PermissionRepo) GetByAppID(appID any) ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.Where("app_id = ?", appID).Find(&permissions).Error
	return permissions, err
}

func (r *PermissionRepo) GetByNames(names []string) ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.Where("name in ?", names).Find(&permissions).Error
	return permissions, err
}
