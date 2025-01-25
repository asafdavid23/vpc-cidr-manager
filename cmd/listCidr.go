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

// listCidrCmd represents the listCidr command
var listCidrCmd = &cobra.Command{
	Use:   "list-cidr",
	Short: "List all CIDRs in the DynamoDB table",
	Run: func(cmd *cobra.Command, args []string) {
		logLevel := viper.GetString("global.logLevel")
		logger := logging.NewLogger(logLevel)
		ctx := context.TODO()
		output := viper.GetString("global.output")
		tableName := viper.GetString("dynamodb.tableName")

		if tableName == "" {
			logger.Fatal("DDB_TABLE_NAME environment variable is not set")
		}

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

		logger.Debug("Listing CIDRs")
		err = internalAws.ListCIDRs(ctx, client, tableName, output)

		if err != nil {
			logger.Fatal(err)
		}
	},
}

func init() {
	// rootCmd.AddCommand(listCidrCmd)
	dynamodbCmd.AddCommand(listCidrCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCidrCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCidrCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
