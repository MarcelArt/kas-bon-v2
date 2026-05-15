package services

import (
	"fmt"

	"github.com/MarcelArt/kas-bon-v2/internal/v1/models"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/repositories"
	"github.com/MarcelArt/kas-bon-v2/pkg/arrays"
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
	GetUsers(id any) ([]models.DomainUser, error)
	GetUserDomains(c fiber.Ctx, id any, parentID any) (page paginate.Page, domains []models.Domain)
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

func (s *DomainService) GetUsers(id any) ([]models.DomainUser, error) {
	res := make([]models.DomainUser, 0)
	domain, err := s.repo.GetByID(id)
	if err != nil {
		return res, err
	}

	policies, err := s.e.GetFilteredGroupingPolicy(0, "", "", domain.Name)
	if err != nil {
		return res, fmt.Errorf("failed to get user lists from policy: %w", err)
	}

	groupedPolicies := make(map[string][][]string)
	groupedPolicies = arrays.Reduce(policies, groupedPolicies, func(total map[string][][]string, currentValue []string) map[string][][]string {
		username := currentValue[0]

		if _, ok := total[username]; !ok {
			total[username] = make([][]string, 0)
		}
		total[username] = append(total[username], currentValue)

		return total
	})

	for k, v := range groupedPolicies {
		user, err := s.uRepo.GetByUsernameOrEmail(k)
		if err != nil {
			continue
		}

		res = append(res, models.DomainUser{
			User:     user,
			Policies: v,
		})
	}

	// for _, policy := range policies {
	// 	user, err := s.uRepo.GetByUsernameOrEmail(policy[0])
	// 	if err != nil {
	// 		continue
	// 	}

	// 	res = append(res, models.DomainUser{
	// 		User:   user,
	// 		Policy: policy,
	// 	})
	// }

	return res, nil
}

func (s *DomainService) GetUserDomains(c fiber.Ctx, id any, parentID any) (page paginate.Page, domains []models.Domain) {
	user, err := s.uRepo.GetByID(id)
	if err != nil {
		page.Error = true
		page.RawError = err
		page.ErrorMessage = fmt.Sprintf("failed retrieving user: %s", err.Error())
		return page, nil
	}

	doms, err := s.e.GetDomainsForUser(user.Username)
	if err != nil {
		page.Error = true
		page.RawError = err
		page.ErrorMessage = fmt.Sprintf("failed policies of user: %s", err.Error())
		return page, nil
	}

	return s.repo.GetDomainByNamesAndParentID(c, parentID, doms)

}
