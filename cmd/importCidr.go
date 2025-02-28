/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"

	internalAws "github.com/asafdavid23/vpc-cidr-manager/internal/aws"
	"github.com/asafdavid23/vpc-cidr-manager/internal/logging"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// importCidrCmd represents the importCidr command
var importCidrCmd = &cobra.Command{
	Use:   "import-cidr",
	Short: "Import CIDR blocks from AWS Env",
	Run: func(cmd *cobra.Command, args []string) {
		logLevel, err := cmd.Flags().GetString("log-level")
		vpcId, err := cmd.Flags().GetString("vpc-id")
		account, err := cmd.Flags().GetString("account-id")
		roleName, err := cmd.Flags().GetString("assume-role")
		tableName := viper.GetString("dynamodb.tableName")
		logger := logging.NewLogger(logLevel)
		ctx := context.TODO()
		assumedRoleArn := "arn:aws:iam::" + account + ":role/earnix/" + roleName
		region := viper.GetString("global.region")

		if region == "" {
			logger.Fatal("AWS_REGION environment variable is not set")
		}

		if tableName == "" {
			logger.Fatal("tableName is not set")
		}

		if vpcId == "" {
			logger.Fatal("vpcId is not set")
		}

		cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))

		if err != nil {
			logger.Fatal(err)
		}

		logger.Debug("Initializing EC2 client")
		hubEC2Client, err := internalAws.GetEc2Client(cfg)

		if err != nil {
			logger.Fatal(err)
		}

		logger.Debug("Initializing DynamoDB client")
		hubDynamoClient, err := internalAws.GetDynamoDBClient(cfg)

		logger.Debug("Initializing STS client")

		if err != nil {
			logger.Fatal(err)
		}

		// Initialize the STS client
		hubStsClient, err := internalAws.GetStsClient(cfg)

		if err != nil {
			logger.Fatal(err)
		}

		if account != "" {
			logger.Debugf("Assuming role for account %s, role %s", account, assumedRoleArn)
			assumedRoleCfg, err := internalAws.AssumeRole(cfg, hubStsClient, assumedRoleArn)

			if err != nil {
				logger.Fatal(err)
			}

			logger.Debugf("%s Role assumed successfully", assumedRoleArn)

			// Extract credentials from the assumed role output
			spokeEC2Client := ec2.NewFromConfig(assumedRoleCfg)

			logger.Debugf("Getting VPC info for vpc %s", vpcId)
			vpcInfo, err := internalAws.GetVpcInfo(spokeEC2Client, vpcId)

			if err != nil {
				logger.Fatal(err)
			}

			logger.Debugf("vpcInfo %v", vpcInfo)

			logger.Debugf("Importing CIDR blocks for vpc %s", vpcId)
			err = internalAws.PushToDynamoDB(ctx, hubDynamoClient, vpcInfo, tableName)

			if err != nil {
				logger.Fatal(err)
			}

			logger.Infof("CIDR block imported successfully")
		}

		logger.Debug("Getting VPC info")
		vpcInfo, err := internalAws.GetVpcInfo(hubEC2Client, vpcId)

		if err != nil {
			logger.Fatal(err)
		}

		logger.Debug("Importing CIDR blocks")
		err = internalAws.PushToDynamoDB(ctx, hubDynamoClient, vpcInfo, tableName)

		if err != nil {
			logger.Fatal(err)
		}

		logger.Infof("CIDR block imported successfully")
	},
}

func init() {
	// rootCmd.AddCommand(importCidrCmd)
	dynamodbCmd.AddCommand(importCidrCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// importCidrCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// importCidrCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	importCidrCmd.Flags().StringP("vpc-id", "v", "", "The VPC ID to import CIDR blocks from")
	importCidrCmd.Flags().StringP("account-id", "a", "", "The AWS account ID to import CIDR blocks from")
	importCidrCmd.Flags().String("assume-role", "", "The role name to assume")
}
