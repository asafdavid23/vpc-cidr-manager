/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"

	internalAws "github.com/asafdavid23/vpc-cidr-manager/internal/aws"
	"github.com/asafdavid23/vpc-cidr-manager/internal/logging"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/spf13/cobra"
)

// importCidrCmd represents the importCidr command
var importCidrCmd = &cobra.Command{
	Use:   "import-cidr",
	Short: "Import CIDR blocks from AWS Env",
	Run: func(cmd *cobra.Command, args []string) {
		logLevel, err := cmd.Flags().GetString("log-level")
		vpcId, err := cmd.Flags().GetString("vpc-id")
		account, err := cmd.Flags().GetString("account-id")
		roleName, err := cmd.Flags().GetString("role-name")

		logger := logging.NewLogger(logLevel)

		logger.Debug("Initializing EC2 client")
		client, err := internalAws.GetEc2Client()

		if err != nil {
			logger.Fatal(err)
		}

		logger.Debug("Initializing DynamoDB client")
		dynamoClient, err := internalAws.GetDynamoDBClient()

		logger.Debug("Initializing STS client")

		if err != nil {
			logger.Fatal(err)
		}

		// Initialize the STS client
		// stsClient, err := internalAws.GetStsClient()

		if err != nil {
			logger.Fatal(err)
		}

		if account != "" {
			logger.Debug("Assuming role")
			assumedRoleOutput, err := internalAws.AssumeRole(roleName, account)

			if err != nil {
				logger.Fatal(err)
			}

			logger.Debug("Role assumed successfully")

			// Extract credentials from the assumed role output
			assumedCredentials := assumedRoleOutput.Credentials

			// Create a new AWS client with the assumed role credentials
			assumedCfg, err := config.LoadDefaultConfig(context.TODO(),
				config.WithCredentialsProvider(
					credentials.NewStaticCredentialsProvider(
						*assumedCredentials.AccessKeyId,
						*assumedCredentials.SecretAccessKey,
						*assumedCredentials.SessionToken,
					),
				),
			)

			if err != nil {
				logger.Fatal(err)
			}

			assumedClient := ec2.NewFromConfig(assumedCfg)

			vpcInfo, err := internalAws.GetVpcInfo(assumedClient, vpcId)

			if err != nil {
				logger.Fatal(err)
			}

			logger.Debug("Importing CIDR blocks")
			err = internalAws.PushToDynamoDB(dynamoClient, vpcInfo)

			if err != nil {
				logger.Fatal(err)
			}

			logger.Infof("CIDR block imported successfully")
		}

		logger.Debug("Getting VPC info")
		vpcInfo, err := internalAws.GetVpcInfo(client, vpcId)

		if err != nil {
			logger.Fatal(err)
		}

		logger.Debug("Importing CIDR blocks")
		err = internalAws.PushToDynamoDB(dynamoClient, vpcInfo)

		if err != nil {
			logger.Fatal(err)
		}

		logger.Infof("CIDR block imported successfully")
	},
}

func init() {
	rootCmd.AddCommand(importCidrCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// importCidrCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// importCidrCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	importCidrCmd.Flags().StringP("vpc-id", "v", "", "The VPC ID to import CIDR blocks from")
	importCidrCmd.MarkFlagRequired("vpc-id")
	importCidrCmd.Flags().StringP("log-level", "l", "info", "The log level to use")
	importCidrCmd.Flags().StringP("account-id", "a", "", "The AWS account ID to import CIDR blocks from")
	importCidrCmd.Flags().String("role-name", "", "The role name to assume")
}
