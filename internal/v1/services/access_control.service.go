package services

import (
	"net/url"

	"github.com/MarcelArt/kas-bon-v2/internal/common"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/models"
	"github.com/casbin/casbin/v3"
)

type IAccessControlService interface {
	GetAllRoles(domain string) ([]string, error)
	GetPermissionsForUser(user string, domain string) ([][]string, error)
	Eval(req models.AccessControlEval) bool
}

type AccessControlService struct {
	e *casbin.Enforcer
}

func NewAccessControlService(e *casbin.Enforcer) *AccessControlService {
	return &AccessControlService{e: e}
}

func (s *AccessControlService) GetAllRoles(domain string) ([]string, error) {
	return s.e.GetAllRolesByDomain(domain)
}

func (s *AccessControlService) GetPermissionsForUser(user string, domain string) ([][]string, error) {
	user, _ = url.QueryUnescape(user)
	domain, _ = url.QueryUnescape(domain)

	return s.e.GetImplicitPermissionsForUser(user, domain)
}

func (s *AccessControlService) Eval(req models.AccessControlEval) bool {
	return common.IsAuthorized(s.e, req.Sub, req.App, req.Dom, req.Obj, req.Act)
}
