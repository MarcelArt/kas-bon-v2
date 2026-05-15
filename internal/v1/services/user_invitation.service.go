package services

import (
	"time"

	"github.com/MarcelArt/kas-bon-v2/internal/v1/models"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/repositories"
	"github.com/casbin/casbin/v3"
	"github.com/gofiber/fiber/v3"
	"github.com/morkid/paginate"
)

type IUserInvitationService interface {
	Create(invitation models.UserInvitationInput) (uint, error)
	Read(c fiber.Ctx) (paginate.Page, []models.UserInvitationPage)
	ReadByUserID(c fiber.Ctx, userID any) (paginate.Page, []models.UserInvitationPage)
	Update(id any, invitation models.UserInvitationInput) error
	Delete(id any) error
	GetByID(id any) (models.UserInvitation, error)
}

type UserInvitationService struct {
	repo  repositories.IUserInvitationRepo
	uRepo repositories.IUserRepo
	dRepo repositories.IDomainRepo
	rRepo repositories.IRoleRepo
	e     *casbin.Enforcer
}

func NewUserInvitationService(repo repositories.IUserInvitationRepo) *UserInvitationService {
	return &UserInvitationService{
		repo: repo,
	}
}

func (s *UserInvitationService) Create(invitation models.UserInvitationInput) (uint, error) {
	return s.repo.Create(invitation)
}

func (s *UserInvitationService) Read(c fiber.Ctx) (paginate.Page, []models.UserInvitationPage) {
	return s.repo.Read(c)
}

func (s *UserInvitationService) ReadByUserID(c fiber.Ctx, userID any) (paginate.Page, []models.UserInvitationPage) {
	return s.repo.ReadByUserID(c, userID)
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

func (s *UserInvitationService) Accept(id any) error {
	today := time.Now()

	invite, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	domain, err := s.dRepo.GetByID(invite.DomainID)
	if err != nil {
		return err
	}

	user, err := s.uRepo.GetByID(invite.UserID)
	if err != nil {
		return err
	}

	role, err := s.repo.GetByID(invite.RoleID)
	if err != nil {
		return err
	}

	s.e.AddGroupingPolicy(user.Username, role.Role.Name, domain.Name)

	s.e.LoadPolicy()
	return s.repo.Update(id, models.UserInvitationInput{AcceptedAt: &today})
}

func (s *UserInvitationService) Reject(id any) error {
	today := time.Now()
	return s.repo.Update(id, models.UserInvitationInput{RejectedAt: &today})
}
