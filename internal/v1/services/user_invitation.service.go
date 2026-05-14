package services

import (
	"github.com/MarcelArt/kas-bon-v2/internal/v1/models"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/repositories"
	"github.com/gofiber/fiber/v3"
	"github.com/morkid/paginate"
)

type IUserInvitationService interface {
	Create(invitation models.UserInvitationInput) (uint, error)
	Read(c fiber.Ctx) (paginate.Page, []models.UserInvitation)
	Update(id any, invitation models.UserInvitationInput) error
	Delete(id any) error
	GetByID(id any) (models.UserInvitation, error)
}

type UserInvitationService struct {
	repo repositories.IUserInvitationRepo
}

func NewUserInvitationService(repo repositories.IUserInvitationRepo) *UserInvitationService {
	return &UserInvitationService{
		repo: repo,
	}
}

func (s *UserInvitationService) Create(invitation models.UserInvitationInput) (uint, error) {
	return s.repo.Create(invitation)
}

func (s *UserInvitationService) Read(c fiber.Ctx) (paginate.Page, []models.UserInvitation) {
	return s.repo.Read(c)
}

func (s *UserInvitationService) Update(id any, invitation models.UserInvitationInput) error {
	return s.repo.Update(id, invitation)
}

func (s *UserInvitationService) Delete(id any) error {
	return s.repo.Delete(id)
}

func (s *UserInvitationService) GetByID(id any) (models.UserInvitation, error) {
	return s.repo.GetByID(id)
}
