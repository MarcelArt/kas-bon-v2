package repositories

import (
	"github.com/MarcelArt/kas-bon-v2/internal/v1/models"
	"github.com/gofiber/fiber/v3"
	"github.com/morkid/paginate"
	"gorm.io/gorm"
)

type IAppRepo interface {
	Create(app models.AppInput) (uint, error)
	Read(c fiber.Ctx) (paginate.Page, []models.App)
	Update(id any, app models.App) error
	Delete(id any) error
	GetByID(id any) (models.App, error)
}

type AppRepo struct {
	db        *gorm.DB
	pageQuery string
}

func NewAppRepo(db *gorm.DB) *AppRepo {
	return &AppRepo{
		db: db,
		pageQuery: `
			select * from apps where deleted_at isnull
		`,
	}
}

func (r *AppRepo) Create(app models.AppInput) (uint, error) {
	err := r.db.Create(&app).Error
	return app.ID, err
}

func (r *AppRepo) Read(c fiber.Ctx) (paginate.Page, []models.App) {
	var apps []models.App
	pg := paginate.New()

	stmt := r.db.Raw(r.pageQuery)

	page := pg.With(stmt).Request(c.Request()).Response(&apps)
	return page, apps
}

func (r *AppRepo) Update(id any, app models.App) error {
	return r.db.Model(&models.App{}).Where("id = ?", id).Updates(&app).Error
}

func (r *AppRepo) Delete(id any) error {
	return r.db.Delete(&models.App{}, id).Error
}

func (r *AppRepo) GetByID(id any) (models.App, error) {
	var app models.App
	err := r.db.Where("id = ?", id).First(&app).Error
	return app, err
}
