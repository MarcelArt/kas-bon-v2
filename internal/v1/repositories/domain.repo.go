package repositories

import (
	"github.com/MarcelArt/kas-bon-v2/internal/v1/models"
	"github.com/gofiber/fiber/v3"
	"github.com/morkid/paginate"
	"gorm.io/gorm"
)

type IDomainRepo interface {
	Create(domain models.DomainInput) (uint, error)
	Read(c fiber.Ctx) (paginate.Page, []models.Domain)
	Update(id any, domain models.Domain) error
	Delete(id any) error
	GetByID(id any) (models.Domain, error)
	GetOrganizationsByNames(name []string) ([]models.Domain, error)
}

type DomainRepo struct {
	db        *gorm.DB
	pageQuery string
}

func NewDomainRepo(db *gorm.DB) *DomainRepo {
	return &DomainRepo{
		db: db,
		pageQuery: `
			select * from domains where deleted_at isnull
		`,
	}
}

func (r *DomainRepo) Create(domain models.DomainInput) (uint, error) {
	err := r.db.Create(&domain).Error
	return domain.ID, err
}

func (r *DomainRepo) Read(c fiber.Ctx) (paginate.Page, []models.Domain) {
	var domains []models.Domain
	pg := paginate.New()

	stmt := r.db.Raw(r.pageQuery)

	page := pg.With(stmt).Request(c.Request()).Response(&domains)
	return page, domains
}

func (r *DomainRepo) Update(id any, domain models.Domain) error {
	return r.db.Model(&models.Domain{}).Where("id = ?", id).Updates(&domain).Error
}

func (r *DomainRepo) Delete(id any) error {
	return r.db.Delete(&models.Domain{}, id).Error
}

func (r *DomainRepo) GetByID(id any) (models.Domain, error) {
	var domain models.Domain
	err := r.db.Where("id = ?", id).First(&domain).Error
	return domain, err
}

func (r *DomainRepo) GetOrganizationsByNames(name []string) ([]models.Domain, error) {
	var domains []models.Domain
	err := r.db.Where("name in ? and is_organization = true", name).Find(&domains).Error
	return domains, err
}
