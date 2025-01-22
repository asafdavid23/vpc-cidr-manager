/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	internalAws "github.com/asafdavid23/vpc-cidr-manager/internal/aws"
	"github.com/asafdavid23/vpc-cidr-manager/internal/logging"
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

		logger := logging.NewLogger(logLevel)

		logger.Debug("Initializing DynamoDB client")
		client, err := internalAws.GetDynamoDBClient()

		if err != nil {
			logger.Fatal(err)
		}

		logger.Debug("Reserving CIDR")
		err = internalAws.ReserveCIDR(client, cidr, vpcID, vpcName, logger)

		if err != nil {
			logger.Fatal(err)
		}

		logger.Infof("CIDR %s reserved successfully", cidr)
	},
}

func init() {
	rootCmd.AddCommand(reserveCidrCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// reserveCidrCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// reserveCidrCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	reserveCidrCmd.Flags().StringP("cidr", "c", "", "The CIDR block to reserve")
	reserveCidrCmd.MarkFlagRequired("cidr")
	reserveCidrCmd.Flags().StringP("log-level", "l", "info", "The log level to use")
	reserveCidrCmd.Flags().StringP("vpc-id", "i", "", "The ID of the VPC to associate with the CIDR block")
	reserveCidrCmd.MarkFlagRequired("vpc-id")
	reserveCidrCmd.Flags().StringP("vpc-name", "n", "", "The name of the VPC to associate with the CIDR block")
	reserveCidrCmd.MarkFlagRequired("vpc-name")
}
