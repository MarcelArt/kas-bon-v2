package services

import (
	"github.com/MarcelArt/kas-bon-v2/internal/v1/models"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/repositories"
	"github.com/gofiber/fiber/v3"
	"github.com/morkid/paginate"
)

type IAppService interface {
	Create(app models.AppInput) (uint, error)
	Read(c fiber.Ctx) (paginate.Page, []models.App)
	Update(id any, app models.App) error
	Delete(id any) error
	GetByID(id any) (models.App, error)
}

type AppService struct {
	repo repositories.IAppRepo
}

func NewAppService(repo repositories.IAppRepo) *AppService {
	return &AppService{repo: repo}
}

func (s *AppService) Create(app models.AppInput) (uint, error) {
	return s.repo.Create(app)
}

func (s *AppService) Read(c fiber.Ctx) (paginate.Page, []models.App) {
	return s.repo.Read(c)
}

func (s *AppService) Update(id any, app models.App) error {
	return s.repo.Update(id, app)
}

func (s *AppService) Delete(id any) error {
	return s.repo.Delete(id)
}

func (s *AppService) GetByID(id any) (models.App, error) {
	return s.repo.GetByID(id)
}
