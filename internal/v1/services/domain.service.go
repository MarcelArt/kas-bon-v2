package services

import (
	"fmt"
	"log"

	"github.com/MarcelArt/kas-bon-v2/internal/v1/models"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/repositories"
	"github.com/casbin/casbin/v3"
	"github.com/gofiber/fiber/v3"
	"github.com/morkid/paginate"
)

type IDomainService interface {
	Create(domain models.DomainInput) (uint, error)
	Read(c fiber.Ctx) (paginate.Page, []models.Domain)
	Update(id any, domain models.Domain) error
	Delete(id any) error
	GetByID(id any) (models.Domain, error)
	GetUsers(id any) ([]models.User, error)
}

type DomainService struct {
	repo  repositories.IDomainRepo
	uRepo repositories.IUserRepo
	e     *casbin.Enforcer
}

func NewDomainService(repo repositories.IDomainRepo, uRepo repositories.IUserRepo, e *casbin.Enforcer) *DomainService {
	return &DomainService{
		repo:  repo,
		uRepo: uRepo,
		e:     e,
	}
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

func (s *DomainService) GetUsers(id any) ([]models.User, error) {
	domain, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	log.Println("e :>> ", s.e)
	usernames, err := s.e.GetAllUsersByDomain(domain.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get user lists from policy")
	}

	return s.uRepo.GetByUsernames(usernames)
}
