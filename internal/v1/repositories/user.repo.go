package repositories

import (
	"github.com/MarcelArt/kas-bon-v2/internal/v1/models"
	"github.com/gofiber/fiber/v3"
	"github.com/morkid/paginate"
	"gorm.io/gorm"
)

type IUserRepo interface {
	Create(user models.UserInput) (uint, error)
	Read(c fiber.Ctx) (paginate.Page, []models.User)
	Update(id any, user models.User) error
	Delete(id any) error
	GetByID(id any) (models.User, error)
	GetByUsernameOrEmail(usernameOrEmail string) (models.User, error)
	GetByUsernames(usernames []string) ([]models.User, error)
}

type UserRepo struct {
	db        *gorm.DB
	pageQuery string
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{
		db: db,
		pageQuery: `
			select * from users where deleted_at isnull
		`,
	}
}

func (r *UserRepo) Create(user models.UserInput) (uint, error) {
	err := r.db.Create(&user).Error
	return user.ID, err
}

func (r *UserRepo) Read(c fiber.Ctx) (paginate.Page, []models.User) {
	var users []models.User
	pg := paginate.New()

	stmt := r.db.Raw(r.pageQuery)

	page := pg.With(stmt).Request(c.Request()).Response(&users)
	return page, users
}

func (r *UserRepo) Update(id any, user models.User) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Updates(&user).Error
}

func (r *UserRepo) Delete(id any) error {
	return r.db.Delete(&models.User{}, id).Error
}

func (r *UserRepo) GetByID(id any) (models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", id).First(&user).Error
	return user, err
}

func (r *UserRepo) GetByUsernameOrEmail(usernameOrEmail string) (models.User, error) {
	var user models.User
	err := r.db.Where("username = ? OR email = ?", usernameOrEmail, usernameOrEmail).First(&user).Error
	return user, err
}

func (r *UserRepo) GetByUsernames(usernames []string) ([]models.User, error) {
	var users []models.User
	err := r.db.Where("username in ?", usernames).Find(&users).Error
	return users, err
}
