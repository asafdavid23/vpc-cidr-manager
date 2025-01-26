/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"os"

	internalAws "github.com/asafdavid23/vpc-cidr-manager/internal/aws"
	"github.com/asafdavid23/vpc-cidr-manager/internal/helpers"
	"github.com/asafdavid23/vpc-cidr-manager/internal/logging"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/spf13/cobra"
)

// createTableCmd represents the createTable command
var createTableCmd = &cobra.Command{
	Use:   "create-table",
	Short: "Create the VpcCidrReservations table in DynamoDB",
	Run: func(cmd *cobra.Command, args []string) {
		logLevel, err := cmd.Flags().GetString("log-level")
		tableName, err := cmd.Flags().GetString("name")
		logger := logging.NewLogger(logLevel)
		ctx := context.TODO()
<<<<<<< Updated upstream
		region := os.Getenv("AWS_REGION")
=======
		stackName := "vpc-cidr-manager-dynamodb-table"
		region := viper.GetString("global.region")
>>>>>>> Stashed changes

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
	},
}

func init() {
	// rootCmd.AddCommand(createTableCmd)
	cfnCmd.AddCommand(createTableCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createTableCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createTableCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	createTableCmd.Flags().StringP("name", "n", "", "The name of the table to create")
	createTableCmd.MarkFlagRequired("name")
}
