package services

import (
	"github.com/MarcelArt/kas-bon-v2/internal/common"
	"github.com/MarcelArt/kas-bon-v2/internal/enums"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/models"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/repositories"
	"github.com/casbin/casbin/v3"
)

type ITokenService interface {
	CheckPermission(userID any, payload models.TokenEndpointRequest) (bool, error)
}

type TokenService struct {
	e     *casbin.Enforcer
	uRepo repositories.IUserRepo
	aRepo repositories.IAppRepo
	dRepo repositories.IDomainRepo
}

func NewTokenService(e *casbin.Enforcer, uRepo repositories.IUserRepo, aRepo repositories.IAppRepo, dRepo repositories.IDomainRepo) *TokenService {
	return &TokenService{e: e, uRepo: uRepo, aRepo: aRepo, dRepo: dRepo}
}

func (s *TokenService) CheckPermission(userID any, payload models.TokenEndpointRequest) (bool, error) {
	user, err := s.uRepo.GetByID(userID)
	if err != nil {
		return false, err
	}

	dom, err := s.dRepo.GetByID(payload.DomainID)
	if err != nil {
		return false, err
	}

	app, err := s.aRepo.GetByID(payload.AppID)
	if err != nil {
		return false, err
	}

	res, act := common.ExtractPermissionResourceAndAction(payload.Permission)

	if ok, _ := s.e.Enforce(user.Username, app.Name, dom.Name, enums.ResourceAll, enums.PermissionFull); ok {
		return true, nil
	}

	if ok, _ := s.e.Enforce(user.Username, app.Name, dom.Name, res, act); ok {
		return true, nil
	}

	return false, nil
}
