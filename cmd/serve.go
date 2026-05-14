/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/MarcelArt/kas-bon-v2/internal/configs"
	apiRoutes "github.com/MarcelArt/kas-bon-v2/internal/v1/routes"
	"github.com/MarcelArt/kas-bon-v2/web/routes"
	"github.com/gofiber/contrib/v3/swaggerui"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/repositories"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/services"
	"github.com/casbin/casbin/v3"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gofiber/fiber/v3"
	html "github.com/gofiber/template/html/v2"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		configs.SetupENV()
		configs.ConnectDB()

		engine := html.New("./web/views", ".html")
		app := fiber.New(fiber.Config{
			Views: engine,
		})
		app.Use(cors.New())

		app.Use(swaggerui.New(swaggerui.Config{
			BasePath: "/",
			FilePath: "./docs/swagger.json",
			Path:     "swagger",
			Title:    "Swagger API Docs",
			CacheAge: 60,
		}))

		app.Get("/public/*", static.New("./public"))

		api := app.Group("/api")
		api.Use(logger.New(logger.Config{
			Format:     "[${time}] ${status} - ${method} ${path} - Query: ${queryParams} - Request: ${body} - Response: ${resBody}\n",
			TimeFormat: "2006-01-02 15:04:05",
			TimeZone:   "Local",
		}))
		apiRoutes.SetupRoutes(api)

		a, _ := gormadapter.NewAdapterByDB(configs.DB)
		e, _ := casbin.NewEnforcer("rbac_model.conf", a)
		userSvc := services.NewUserService(
			repositories.NewUserRepo(configs.DB),
			repositories.NewDomainRepo(configs.DB),
			repositories.NewRoleRepo(configs.DB),
			e,
		)
		routes.SetupWebRoutes(app, userSvc, e)

		port := fmt.Sprintf(":%s", configs.Env.PORT)
		log.Printf("Listening on port %s", configs.Env.PORT)
		app.Listen(port)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
