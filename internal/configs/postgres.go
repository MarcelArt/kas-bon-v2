package configs

import (
	"fmt"
	"strconv"

	"github.com/MarcelArt/kas-bon-v2/internal/enums"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB
var dsn string

func ConnectDB() {
	p := Env.DBPort
	port, err := strconv.ParseUint(p, 10, 32)
	dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s search_path=%s sslmode=disable", Env.DBHost, port, Env.DBUser, Env.DBPassword, Env.DBName, Env.DBSchema)

	if err != nil {
		panic("failed to parse database port")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		panic("failed to connect database")
	}

	DB = db

	fmt.Println("Connection Opened to Database")
}

func MigrateDB() error {
	db := DB
	err := db.AutoMigrate(
		models.User{},
		models.Domain{},
		models.App{},
		models.Role{},
		models.Permission{},
		models.UserInvitation{},
	)
	fmt.Println("Database Migrated")

	seedApp()
	seedPermissions()

	return err
}

func DropDB() error {
	db := DB
	err := db.Migrator().DropTable(
		models.User{},
		models.Domain{},
		models.App{},
		models.Role{},
		models.Permission{},
	)
	fmt.Println("Database Dropped")

	return err
}

func seedApp() {
	DB.Create(&models.App{
		Name: enums.AppName,
	})
}

func seedPermissions() {
	permissions := []models.Permission{
		{Name: "all#fullAccess", Description: "Super user", AppID: enums.AppID},
		{Name: "users#read", Description: "Read users", AppID: enums.AppID},
		{Name: "users#update", Description: "Update users", AppID: enums.AppID},
		{Name: "users#delete", Description: "Delete users", AppID: enums.AppID},
		{Name: "roles#read", Description: "Read roles", AppID: enums.AppID},
		{Name: "roles#create", Description: "Create roles", AppID: enums.AppID},
		{Name: "roles#update", Description: "Update roles", AppID: enums.AppID},
		{Name: "roles#delete", Description: "Delete roles", AppID: enums.AppID},
		{Name: "permissions#read", Description: "Read permissions", AppID: enums.AppID},
		{Name: "permissions#create", Description: "Create permissions", AppID: enums.AppID},
		{Name: "permissions#update", Description: "Update permissions", AppID: enums.AppID},
		{Name: "permissions#delete", Description: "Delete permissions", AppID: enums.AppID},
		{Name: "apps#read", Description: "Read apps", AppID: enums.AppID},
		{Name: "apps#create", Description: "Create apps", AppID: enums.AppID},
		{Name: "apps#update", Description: "Update apps", AppID: enums.AppID},
		{Name: "apps#delete", Description: "Delete apps", AppID: enums.AppID},
		{Name: "domains#read", Description: "Read domains", AppID: enums.AppID},
		{Name: "domains#create", Description: "Create domains", AppID: enums.AppID},
		{Name: "domains#update", Description: "Update domains", AppID: enums.AppID},
		{Name: "domains#delete", Description: "Delete domains", AppID: enums.AppID},
		{Name: "user-invitations#read", Description: "Read user invitations", AppID: enums.AppID},
		{Name: "user-invitations#create", Description: "Create user invitations", AppID: enums.AppID},
		{Name: "user-invitations#update", Description: "Update user invitations", AppID: enums.AppID},
		{Name: "user-invitations#delete", Description: "Delete user invitations", AppID: enums.AppID},
	}

	for _, p := range permissions {
		DB.FirstOrCreate(&p, models.Permission{Name: p.Name})
	}
}
