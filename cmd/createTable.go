/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"

	internalAws "github.com/asafdavid23/vpc-cidr-manager/internal/aws"
	"github.com/asafdavid23/vpc-cidr-manager/internal/logging"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// createTableCmd represents the createTable command
var createTableCmd = &cobra.Command{
	Use:   "create-table",
	Short: "Create the VpcCidrReservations table in DynamoDB",
	Run: func(cmd *cobra.Command, args []string) {
		logLevel, err := cmd.Flags().GetString("log-level")
		tableName := viper.GetString("dynamodb.tableName")
		logger := logging.NewLogger(logLevel)
		ctx := context.TODO()
		region := viper.GetString("global.region")

		if region == "" {
			logger.Fatal("AWS_REGION environment variable is not set")
		}

		cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))

		if err != nil {
			logger.Fatal(err)
		}

		logger.Debug("Initializing DynamoDB client")
		client, err := internalAws.GetDynamoDBClient(cfg)
		if err != nil {
			logger.Fatal(err)
		}

		logger.Debug("Creating DynamoDB table")
		err = internalAws.CreateDynamoDBTable(ctx, client, tableName, logger)

		if err != nil {
			logger.Fatal(err)
		}

		logger.Infof("%s DynamoDB table created successfully", tableName)
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

	viper.BindPFlag("dynamodb.tableName", createTableCmd.Flags().Lookup("name"))
}
