/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/MarcelArt/kas-bon-v2/internal/configs"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/routes"
	"github.com/gofiber/contrib/v3/swaggerui"
	"github.com/gofiber/fiber/v3"
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
		// configs.SetupCasbin()

		app := fiber.New()
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
		routes.SetupRoutes(api)

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
