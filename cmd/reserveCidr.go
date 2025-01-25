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

// reserveCidrCmd represents the reserveCidr command
var reserveCidrCmd = &cobra.Command{
	Use:   "reserve-cidr",
	Short: "Reserve a CIDR block",
	Run: func(cmd *cobra.Command, args []string) {
		logLevel, err := cmd.Flags().GetString("log-level")
		cidr, err := cmd.Flags().GetString("cidr")
		vpcID, err := cmd.Flags().GetString("vpc-id")
		vpcName, err := cmd.Flags().GetString("vpc-name")
		ctx := context.TODO()
		logger := logging.NewLogger(logLevel)
		autoGenerate, err := cmd.Flags().GetBool("auto-generate")

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

		if autoGenerate {
			baseCidr, err := cmd.Flags().GetString("base-cidr")
			prefixSize, err := cmd.Flags().GetInt("prefix-size")

			if baseCidr == "" && prefixSize == 0 {
				logger.Fatal("base-cidr and prefix-size flags are required when auto-generate flag is set")
			}

			existingCidrs, err := internalAws.FetchExistingCIDRs(client, tableName)

			logger.Debugf("Fetching existing CIDRs from DynamoDB table %v", existingCidrs)

			if err != nil {
				logger.Fatal(err)
			}

			logger.Debug("Generating CIDR")
			cidr, err = helpers.GenerateCIDR(existingCidrs, baseCidr, prefixSize)

			if err != nil {
				logger.Fatal(err)
			}

			logger.Infof("Generated CIDR: %s", cidr)

		}

		logger.Debug("Reserving CIDR")
		err = internalAws.ReserveCIDR(ctx, client, cidr, vpcID, vpcName, logger)

		if err != nil {
			logger.Fatal(err)
		}

		logger.Infof("CIDR %s reserved successfully", cidr)
	},
}

func init() {
	// rootCmd.AddCommand(reserveCidrCmd)
	dynamodbCmd.AddCommand(reserveCidrCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// reserveCidrCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// reserveCidrCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	reserveCidrCmd.Flags().StringP("cidr", "c", "", "The CIDR block to reserve")
	reserveCidrCmd.Flags().String("log-level", "info", "The log level to use")
	reserveCidrCmd.Flags().String("vpc-id", "", "The ID of the VPC to associate with the CIDR block")
	reserveCidrCmd.Flags().String("vpc-name", "", "The name of the VPC to associate with the CIDR block")
	reserveCidrCmd.Flags().Bool("auto-generate", false, "Automatically generate a CIDR block")
	reserveCidrCmd.Flags().String("base-cidr", "", "The base CIDR block to use when auto-generating a CIDR block")
	reserveCidrCmd.Flags().Int("prefix-size", 16, "The prefix size to use when auto-generating a CIDR block")
}
