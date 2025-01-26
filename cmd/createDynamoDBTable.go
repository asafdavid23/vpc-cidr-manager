/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"

	internalAws "github.com/asafdavid23/vpc-cidr-manager/internal/aws"
	"github.com/asafdavid23/vpc-cidr-manager/internal/helpers"
	"github.com/asafdavid23/vpc-cidr-manager/internal/logging"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// createTableCmd represents the createTable command
var createDynamoDBTableCmd = &cobra.Command{
	Use:   "dynamodb-table",
	Short: "Create the VpcCidrReservations table in DynamoDB",
	Run: func(cmd *cobra.Command, args []string) {
		logLevel, err := cmd.Flags().GetString("log-level")
		tableName := viper.GetString("dynamodb.tableName")
		logger := logging.NewLogger(logLevel)
		ctx := context.TODO()
		stackName := "vpc-cidr-manager-dynamodb-table"
		dryRun, err := cmd.Flags().GetBool("dry-run")
		generateTemplate, err := cmd.Flags().GetBool("generate-iaac-template")
		region := viper.GetString("global.region")

		if region == "" {
			logger.Fatal("region is not set")
		}

		cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))

		if err != nil {
			logger.Fatal(err)
		}

		dynamodbTableTemplateFile := "templates/cloudformation/dynamodb_table.yml"

		data := helpers.DynamoDBTableTemplateData{
			TableName: tableName,
		}

		renderedTemplate, err := helpers.LoadAndRenderCFNTemplate(dynamodbTableTemplateFile, data)

		if err != nil {
			logger.Fatal(err)
		}

		logger.Debugf("Rendered template: %s", renderedTemplate)

		logger.Debug("Iinitializing CloudFormation client")
		cfnClient, err := internalAws.InitializeCFNClient(cfg)

		if err != nil {
			logger.Fatal(err)
		}

		if dryRun {
			logger.Infof("Dry run enabled, not creating stack %s", stackName)
			err := internalAws.ValidateCFNStackTemplate(ctx, cfnClient, renderedTemplate)

			if err != nil {
				logger.Fatal(err)
			}

		} else {
			logger.Debug("Creating cloudformation stack")

			output, err := internalAws.CreateCFNStack(ctx, cfnClient, stackName, renderedTemplate)

			if err != nil {
				logger.Fatal(err)
			}

			err = internalAws.WaitForStackToBeCreated(ctx, cfnClient, stackName)

			if err != nil {
				logger.Fatal(err)
			}

			logger.Infof("Stack %s created successfully", *output.StackId)
		}
	},
}

func init() {
	// rootCmd.AddCommand(createTableCmd)
	createCmd.AddCommand(createDynamoDBTableCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createTableCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createTableCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	createDynamoDBTableCmd.Flags().StringP("name", "n", "", "The name of the table to create")
	createDynamoDBTableCmd.Flags().Bool("dry-run", false, "Print the rendered CloudFormation template without creating the stack")
	createDynamoDBTableCmd.Flags().Bool("generate-iaac-template", false, "Generate the CloudFormation template without creating the stack")

	viper.BindPFlag("dynamodb.tableName", createDynamoDBTableCmd.Flags().Lookup("name"))
}
