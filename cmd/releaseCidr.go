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

var cidr []string

// releaseCidrCmd represents the releaseCidr command
var releaseCidrCmd = &cobra.Command{
	Use:   "release-cidr",
	Short: "Release a CIDR block",
	Run: func(cmd *cobra.Command, args []string) {
		logLevel, err := cmd.Flags().GetString("log-level")
		logger := logging.NewLogger(logLevel)
		ctx := context.TODO()
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

		logger.Debug("Releasing CIDR block")
		err = internalAws.ReleaseCidr(ctx, client, cidr, logger)

		if err != nil {
			logger.Fatal(err)
		}

		logger.Infof("%s CIDR block released successfully", cidr)
	},
}

func init() {
	// rootCmd.AddCommand(releaseCidrCmd)
	dynamodbCmd.AddCommand(releaseCidrCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// releaseCidrCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// releaseCidrCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	releaseCidrCmd.Flags().StringSliceVarP(&cidr, "cidr", "c", []string{}, "The CIDR block to release")
	releaseCidrCmd.MarkFlagRequired("cidr")
}
