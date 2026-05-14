package services

import (
	"fmt"

	"github.com/MarcelArt/kas-bon-v2/internal/common"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/models"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/repositories"
	"github.com/MarcelArt/kas-bon-v2/pkg/arrays"
	"github.com/casbin/casbin/v3"
	"github.com/gofiber/fiber/v3"
	"github.com/morkid/paginate"
)

type IRoleService interface {
	Create(role models.RoleInput) (uint, error)
	Read(c fiber.Ctx, domainID any) (paginate.Page, []models.Role)
	Update(id any, role models.Role) error
	Delete(id any) error
	GetByID(id any) (models.Role, error)
	GetPermissions(id any) (models.GetRolePermissionsResponse, error)
	AssignPermissions(roleID any, appID any, permissionIDs []uint) ([]string, error)
	GetByDomainID(domainID any) ([]models.Role, error)
}

type RoleService struct {
	repo  repositories.IRoleRepo
	aRepo repositories.IAppRepo
	dRepo repositories.IDomainRepo
	pRepo repositories.IPermissionRepo
	e     *casbin.Enforcer
}

func NewRoleService(repo repositories.IRoleRepo, aRepo repositories.IAppRepo, dRepo repositories.IDomainRepo, pRepo repositories.IPermissionRepo, e *casbin.Enforcer) *RoleService {
	return &RoleService{repo: repo, aRepo: aRepo, dRepo: dRepo, pRepo: pRepo, e: e}
}

func (s *RoleService) Create(role models.RoleInput) (uint, error) {
	return s.repo.Create(role)
}

func (s *RoleService) Read(c fiber.Ctx, domainID any) (paginate.Page, []models.Role) {
	return s.repo.Read(c, domainID)
}

func (s *RoleService) Update(id any, role models.Role) error {
	return s.repo.Update(id, role)
}

func (s *RoleService) Delete(id any) error {
	return s.repo.Delete(id)
}

func (s *RoleService) GetByID(id any) (models.Role, error) {
	return s.repo.GetByID(id)
}

func (s *RoleService) GetPermissions(id any) (models.GetRolePermissionsResponse, error) {
	var response models.GetRolePermissionsResponse
	role, err := s.repo.GetByID(id)
	if err != nil {
		return response, err
	}

	dom, err := s.dRepo.GetByID(role.DomainID)
	if err != nil {
		return response, err
	}

	response.Policies, err = s.e.GetImplicitPermissionsForUser(role.Name, dom.Name)
	if err != nil {
		return response, err
	}

	permKeys := arrays.Map(response.Policies, func(policy []string) string {
		return fmt.Sprintf("%s#%s", policy[3], policy[4])
	})

	response.Permissions, err = s.pRepo.GetByNames(permKeys)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (s *RoleService) AssignPermissions(roleID any, appID any, permissionIDs []uint) ([]string, error) {
	app, err := s.aRepo.GetByID(appID)
	if err != nil {
		return nil, err
	}

	role, err := s.repo.GetByID(roleID)
	if err != nil {
		return nil, err
	}

	dom, err := s.dRepo.GetByID(role.DomainID)
	if err != nil {
		return nil, err
	}

	permissions, err := s.pRepo.GetDistinctByIDs(permissionIDs)
	if err != nil {
		return nil, err
	}

	if _, err := s.e.RemoveFilteredPolicy(0, role.Name, app.Name, dom.Name); err != nil {
		return nil, err
	}

	for _, permission := range permissions {
		res, act := common.ExtractPermissionResourceAndAction(permission)

		if _, err := s.e.AddPolicy(role.Name, app.Name, dom.Name, res, act); err != nil {
			return nil, err
		}
	}

	return permissions, nil
}

func (s *RoleService) GetByDomainID(domainID any) ([]models.Role, error) {
	return s.repo.GetByDomainID(domainID)
}
