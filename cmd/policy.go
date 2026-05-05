/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/MarcelArt/kas-bon-v2/internal/configs"
	"github.com/casbin/casbin/v3"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/spf13/cobra"
)

var user string
var role string
var obj string
var act string

// policyCmd represents the policy command
var policyCmd = &cobra.Command{
	Use:   "policy",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		configs.SetupENV()
		configs.ConnectDB()
		a, _ := gormadapter.NewAdapterByDB(configs.DB)

		e, _ := casbin.NewEnforcer("rbac_model.conf", a)
		e.LoadPolicy()

		e.AddPolicy(role, obj, act)
		e.AddGroupingPolicy(user, role)

		fmt.Println("Policy created", role, obj, act)
	},
}

func init() {
	rootCmd.AddCommand(policyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// policyCmd.PersistentFlags().String("foo", "", "A help for foo")
	policyCmd.PersistentFlags().StringVarP(&user, "user", "u", "kandar", "user")
	policyCmd.PersistentFlags().StringVarP(&role, "role", "r", "viewer", "role")
	policyCmd.PersistentFlags().StringVarP(&obj, "obj", "o", "users", "resource/object")
	policyCmd.PersistentFlags().StringVarP(&act, "act", "a", "read", "action")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// policyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
