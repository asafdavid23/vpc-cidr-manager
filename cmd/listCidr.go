/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"os"

	internalAws "github.com/asafdavid23/vpc-cidr-manager/internal/aws"
	"github.com/asafdavid23/vpc-cidr-manager/internal/logging"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/spf13/cobra"
)

// listCidrCmd represents the listCidr command
var listCidrCmd = &cobra.Command{
	Use:   "list-cidr",
	Short: "List all CIDRs in the DynamoDB table",
	Run: func(cmd *cobra.Command, args []string) {
		logLevel, err := cmd.Flags().GetString("log-level")
		logger := logging.NewLogger(logLevel)
		ctx := context.TODO()
		output, err := cmd.Flags().GetString("output")

		tableName := os.Getenv("DDB_TABLE_NAME")

		if tableName == "" {
			logger.Fatal("DDB_TABLE_NAME environment variable is not set")
		}

		region := os.Getenv("AWS_REGION")

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
	rootCmd.AddCommand(listCidrCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCidrCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCidrCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	listCidrCmd.Flags().StringP("output", "o", "table", "Output format (table, json, yaml)")
	listCidrCmd.Flags().StringP("log-level", "l", "info", "Log level (debug, info, warn, error, fatal)")
}
