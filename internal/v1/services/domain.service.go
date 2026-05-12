package services

import (
	"github.com/MarcelArt/kas-bon-v2/internal/v1/models"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/repositories"
	"github.com/gofiber/fiber/v3"
	"github.com/morkid/paginate"
)

type IDomainService interface {
	Create(domain models.DomainInput) (uint, error)
	Read(c fiber.Ctx) (paginate.Page, []models.Domain)
	Update(id any, domain models.Domain) error
	Delete(id any) error
	GetByID(id any) (models.Domain, error)
}

type DomainService struct {
	repo repositories.IDomainRepo
}

func NewDomainService(repo repositories.IDomainRepo) *DomainService {
	return &DomainService{repo: repo}
}

func (s *DomainService) Create(domain models.DomainInput) (uint, error) {
	return s.repo.Create(domain)
}

func (s *DomainService) Read(c fiber.Ctx) (paginate.Page, []models.Domain) {
	return s.repo.Read(c)
}

func (s *DomainService) Update(id any, domain models.Domain) error {
	return s.repo.Update(id, domain)
}

func (s *DomainService) Delete(id any) error {
	return s.repo.Delete(id)
}

func (s *DomainService) GetByID(id any) (models.Domain, error) {
	return s.repo.GetByID(id)
}
