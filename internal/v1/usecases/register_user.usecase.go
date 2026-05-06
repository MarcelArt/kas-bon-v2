package usecases

import (
	"fmt"

	"github.com/MarcelArt/kas-bon-v2/internal/enums"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/models"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/repositories"
	"github.com/alexedwards/argon2id"
	"github.com/casbin/casbin/v3"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

type RegisterUserUsecase struct {
	User models.UserInput

	uRepo repositories.IUserRepo
	dRepo repositories.IDomainRepo
	e     *casbin.Enforcer
}

func InitRegisterUserUsecase(tx *gorm.DB) *RegisterUserUsecase {
	a, _ := gormadapter.NewAdapterByDB(tx)

	e, _ := casbin.NewEnforcer("rbac_model.conf", a)

	return &RegisterUserUsecase{
		uRepo: repositories.NewUserRepo(tx),
		dRepo: repositories.NewDomainRepo(tx),
		e:     e,
	}
}

func (u *RegisterUserUsecase) Execute() (uint, error) {
	user := u.User

	dom := fmt.Sprintf("%s's organization", u.User.Username)
	role := "Owner"

	u.e.AddPolicy(role, enums.AppName, dom, enums.ResourceAll, enums.PermissionFull)
	u.e.AddGroupingPolicy(u.User.Username, role, dom)

	password, err := argon2id.CreateHash(user.Password, argon2id.DefaultParams)
	if err != nil {
		return 0, fmt.Errorf("failed hashing password: %w", err)
	}
	user.Password = password

	_, err = u.dRepo.Create(models.DomainInput{
		Name:           dom,
		IsOrganization: true,
	})
	if err != nil {
		return 0, fmt.Errorf("failed creating domain: %w", err)
	}

	return u.uRepo.Create(user)
}
