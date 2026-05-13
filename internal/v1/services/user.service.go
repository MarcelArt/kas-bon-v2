package services

import (
	"fmt"

	"github.com/MarcelArt/kas-bon-v2/internal/common"
	"github.com/MarcelArt/kas-bon-v2/internal/configs"
	"github.com/MarcelArt/kas-bon-v2/internal/enums"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/models"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/repositories"
	"github.com/alexedwards/argon2id"
	"github.com/casbin/casbin/v3"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gofiber/fiber/v3"
	"github.com/morkid/paginate"
)

type IUserService interface {
	Create(user models.UserInput) (uint, error)
	Read(c fiber.Ctx) (paginate.Page, []models.User)
	Update(id any, user models.User) error
	Delete(id any) error
	GetByID(id any) (models.User, error)
	Login(login models.LoginInput, c fiber.Ctx) (models.LoginResponse, error)
	Refresh(userID any, isRemember bool, c fiber.Ctx) (models.LoginResponse, error)
	GetRoles(id any, domainID any) ([]string, error)
	GetPermissions(id any, domainID any) ([][]string, error)
	AssignRoles(id any, domainID any, roleIDs []uint) ([]string, error)
	GetOrganizations(id any) ([]models.Domain, error)
}

type UserService struct {
	repo  repositories.IUserRepo
	dRepo repositories.IDomainRepo
	rRepo repositories.IRoleRepo
	e     *casbin.Enforcer
}

func NewUserService(repo repositories.IUserRepo, dRepo repositories.IDomainRepo, rRepo repositories.IRoleRepo, e *casbin.Enforcer) *UserService {
	return &UserService{repo: repo, dRepo: dRepo, rRepo: rRepo, e: e}
}

func (s *UserService) Create(user models.UserInput) (uint, error) {
	tx := configs.DB.Begin()
	defer tx.Rollback()

	a, _ := gormadapter.NewAdapterByDB(tx)
	enforcer, _ := casbin.NewEnforcer("rbac_model.conf", a)

	dom := fmt.Sprintf("%s's organization", user.Username)

	enforcer.AddPolicy(enums.RoleDefault, enums.AppName, dom, enums.ResourceAll, enums.PermissionFull)
	enforcer.AddGroupingPolicy(user.Username, enums.RoleDefault, dom)

	password, err := argon2id.CreateHash(user.Password, argon2id.DefaultParams)
	if err != nil {
		return 0, fmt.Errorf("failed hashing password: %w", err)
	}
	user.Password = password

	uRepo := repositories.NewUserRepo(tx)
	dRepoTx := repositories.NewDomainRepo(tx)
	rRepoTx := repositories.NewRoleRepo(tx)
	pRepoTx := repositories.NewPermissionRepo(tx)

	domainID, err := dRepoTx.Create(models.DomainInput{
		Name:           dom,
		IsOrganization: true,
	})
	if err != nil {
		return 0, fmt.Errorf("failed creating domain: %w", err)
	}

	if _, err := rRepoTx.Create(models.RoleInput{Name: enums.RoleDefault, DomainID: domainID}); err != nil {
		return 0, fmt.Errorf("failed creating default role: %w", err)
	}

	permission := fmt.Sprintf("%s#%s", enums.ResourceAll, enums.PermissionFull)
	if _, err := pRepoTx.Create(models.PermissionInput{Name: permission, AppID: enums.AppID}); err != nil {
		return 0, fmt.Errorf("failed creating default permission: %w", err)
	}

	id, err := uRepo.Create(user)
	if err != nil {
		return 0, err
	}

	tx.Commit()
	return id, nil
}

func (s *UserService) Read(c fiber.Ctx) (paginate.Page, []models.User) {
	return s.repo.Read(c)
}

func (s *UserService) Update(id any, user models.User) error {
	return s.repo.Update(id, user)
}

func (s *UserService) Delete(id any) error {
	return s.repo.Delete(id)
}

func (s *UserService) GetByID(id any) (models.User, error) {
	return s.repo.GetByID(id)
}

func (s *UserService) Login(login models.LoginInput, c fiber.Ctx) (models.LoginResponse, error) {
	user, err := s.repo.GetByUsernameOrEmail(login.Username)
	if err != nil {
		return models.LoginResponse{}, err
	}

	ok, err := argon2id.ComparePasswordAndHash(login.Password, user.Password)
	if err != nil {
		return models.LoginResponse{}, err
	}

	if !ok {
		return models.LoginResponse{}, fmt.Errorf("invalid credentials")
	}

	return s.generateTokenPair(user, login.IsRemember, c)
}

func (s *UserService) Refresh(userID any, isRemember bool, c fiber.Ctx) (models.LoginResponse, error) {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return models.LoginResponse{}, err
	}

	return s.generateTokenPair(user, isRemember, c)
}

func (s *UserService) generateTokenPair(user models.User, isRemember bool, c fiber.Ctx) (models.LoginResponse, error) {
	a, _ := gormadapter.NewAdapterByDB(configs.DB)
	e, _ := casbin.NewEnforcer("rbac_model.conf", a)

	permissions, err := e.GetImplicitPermissionsForUser(user.Username)
	if err != nil {
		return models.LoginResponse{}, fmt.Errorf("failed retrieving permissions: %w", err)
	}

	claims := map[string]any{
		"sub":    user.Username,
		"userId": user.ID,
		"iss":    c.BaseURL(),
	}
	at, rt, err := common.GenerateJWTPair(claims, permissions, isRemember)
	if err != nil {
		return models.LoginResponse{}, fmt.Errorf("failed generating tokens: %w", err)
	}

	return models.LoginResponse{
		AccessToken:  at,
		RefreshToken: rt,
		User:         user,
	}, nil
}

func (s *UserService) GetRoles(id any, domainID any) ([]string, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	domain, err := s.dRepo.GetByID(domainID)
	if err != nil {
		return nil, err
	}

	roles := s.e.GetRolesForUserInDomain(user.Username, domain.Name)
	return roles, nil
}

func (s *UserService) GetPermissions(id any, domainID any) ([][]string, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	dom, err := s.dRepo.GetByID(domainID)
	if err != nil {
		return nil, err
	}

	permissions, err := s.e.GetImplicitPermissionsForUser(user.Username, dom.Name)
	if err != nil {
		return nil, err
	}

	return permissions, nil
}

func (s *UserService) AssignRoles(id any, domainID any, roleIDs []uint) ([]string, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	dom, err := s.dRepo.GetByID(domainID)
	if err != nil {
		return nil, err
	}

	if _, err := s.e.RemoveFilteredGroupingPolicy(0, user.Username, "", dom.Name); err != nil {
		return nil, err
	}

	roles, err := s.rRepo.GetDistinctByIDs(roleIDs)
	if err != nil {
		return nil, err
	}

	for _, role := range roles {
		if _, err := s.e.AddGroupingPolicy(user.Username, role, dom.Name); err != nil {
			return nil, err
		}
	}

	return roles, nil
}

func (s *UserService) GetOrganizations(id any) ([]models.Domain, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	doms, err := s.e.GetDomainsForUser(user.Username)
	if err != nil {
		return nil, err
	}

	domains, err := s.dRepo.GetOrganizationsByNames(doms)
	if err != nil {
		return nil, err
	}

	return domains, nil
}
