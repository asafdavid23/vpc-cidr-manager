/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/asafdavid23/vpc-cidr-manager/internal/helpers"
	"github.com/asafdavid23/vpc-cidr-manager/internal/logging"
	"github.com/spf13/cobra"
)

// createAssumedRoleCmd represents the createAssumedRule command
var createAssumedRoleCmd = &cobra.Command{
	Use:   "create-assumed-role",
	Short: "Create an assumed role for the VPC CIDR Manager",
	Run: func(cmd *cobra.Command, args []string) {
		logLevel, err := cmd.Flags().GetString("log-level")
		roleName, err := cmd.Flags().GetString("role-name")
		logger := logging.NewLogger(logLevel)
		assumeRolePrincipal, err := cmd.Flags().GetString("assume-role-principal")

		iamTemplateFile := "templates/cloudformation/iam_role.yml"

		data := helpers.IAMTemplateData{
			RoleName:  roleName,
			Principal: assumeRolePrincipal,
		}

		renderedTemplate, err := helpers.LoadAndRenderTemplate(iamTemplateFile, data)

		if err != nil {
			logger.Fatal(err)
		}

		// Print the rendered CloudFormation template
		fmt.Println("Rendered CloudFormation Template:")

		// logger.Debug("Initializing IAM client")
		// client, err := internalAws.GetIAMClient()

		// if err != nil {
		// 	logger.Fatal(err)
		// }

		// logger.Debug("Creating assumed role")
		// err = internalAws.CreateAssumableRole(client, roleName, policyFile, trustFile)

		// if err != nil {
		// 	logger.Fatal(err)
		// }

		// logger.Infof("Assumed role %s created successfully", roleName)
	},
}

func init() {
	rootCmd.AddCommand(createAssumedRoleCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createAssumedRuleCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createAssumedRuleCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	createAssumedRoleCmd.Flags().StringP("log-level", "l", "info", "The log level to use")
	createAssumedRoleCmd.Flags().StringP("role-name", "r", "", "The name of the role to create")
	createAssumedRoleCmd.MarkFlagRequired("role-name")
	createAssumedRoleCmd.Flags().String("assume-role-principal", "", "The principal to assume the role")

}
