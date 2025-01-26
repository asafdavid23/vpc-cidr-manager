/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"

	"github.com/asafdavid23/vpc-cidr-manager/internal/aws"
	internalAws "github.com/asafdavid23/vpc-cidr-manager/internal/aws"
	"github.com/asafdavid23/vpc-cidr-manager/internal/helpers"
	"github.com/asafdavid23/vpc-cidr-manager/internal/logging"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// createAssumedRoleCmd represents the createAssumedRule command
var createAssumedRoleCmd = &cobra.Command{
	Use:   "assumed-role",
	Short: "Create an assumed role for the VPC CIDR Manager",
	Run: func(cmd *cobra.Command, args []string) {
		logLevel, err := cmd.Flags().GetString("log-level")
		roleName, err := cmd.Flags().GetString("role-name")
		logger := logging.NewLogger(logLevel)
		hubAccount, err := cmd.Flags().GetString("hub-account")
		assumeRolePrincipal := "arn:aws:iam::" + hubAccount + ":root"
		ctx := context.TODO()
		stackName := "vpc-cidr-manager-assumed-role"

		region := viper.GetString("global.region")

		if region == "" {
			logger.Fatal("region is not set")
		}

		cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))

		if err != nil {
			logger.Fatal(err)
		}

		iamTemplateFile := "templates/cloudformation/iam_role.yml"

		data := helpers.IAMTemplateData{
			RoleName:  roleName,
			Principal: assumeRolePrincipal,
		}

		renderedTemplate, err := helpers.LoadAndRenderIAMTemplate(iamTemplateFile, data)

		if err != nil {
			logger.Fatal(err)
		}

		logger.Debugf("Rendered template: %s", renderedTemplate)

		logger.Debug("Iinitializing CloudFormation client")
		cfnClient, err := aws.InitializeCFNClient(cfg)

		if err != nil {
			logger.Fatal(err)
		}

		logger.Debug("Creating CloudFormation stack")

		output, err := internalAws.CreateCFNStack(ctx, cfnClient, stackName, renderedTemplate)

		if err != nil {
			logger.Fatal(err)
		}

		err = internalAws.WaitForStackToBeCreated(ctx, cfnClient, stackName)

		if err != nil {
			logger.Fatal(err)
		}

		logger.Infof("Stack %s created successfully", *output.StackId)
	},
}

func init() {
	// rootCmd.AddCommand(createAssumedRoleCmd)
	createCmd.AddCommand(createAssumedRoleCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createAssumedRuleCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createAssumedRuleCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	createAssumedRoleCmd.Flags().StringP("role-name", "r", "", "The name of the role to create")
	createAssumedRoleCmd.Flags().String("hub-account", "", "The principal that assume the role in spoke accoount")

	viper.BindPFlag("iam.hubAccountId", createAssumedRoleCmd.Flags().Lookup("hub-account"))
	viper.BindPFlag("iam.assumedRoleName", createAssumedRoleCmd.Flags().Lookup("role-name"))
}
