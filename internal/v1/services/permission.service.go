package services

import (
	"github.com/MarcelArt/kas-bon-v2/internal/v1/models"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/repositories"
	"github.com/gofiber/fiber/v3"
	"github.com/morkid/paginate"
)

type IPermissionService interface {
	Create(permission models.PermissionInput) (uint, error)
	Read(c fiber.Ctx, appID any) (paginate.Page, []models.Permission)
	Update(id any, permission models.Permission) error
	Delete(id any) error
	GetByID(id any) (models.Permission, error)
}

type PermissionService struct {
	repo repositories.IPermissionRepo
}

func NewPermissionService(repo repositories.IPermissionRepo) *PermissionService {
	return &PermissionService{repo: repo}
}

func (s *PermissionService) Create(permission models.PermissionInput) (uint, error) {
	return s.repo.Create(permission)
}

func (s *PermissionService) Read(c fiber.Ctx, appID any) (paginate.Page, []models.Permission) {
	return s.repo.Read(c, appID)
}

func (s *PermissionService) Update(id any, permission models.Permission) error {
	return s.repo.Update(id, permission)
}

func (s *PermissionService) Delete(id any) error {
	return s.repo.Delete(id)
}

func (s *PermissionService) GetByID(id any) (models.Permission, error) {
	return s.repo.GetByID(id)
}
