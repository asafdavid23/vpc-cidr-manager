/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	internalAws "github.com/asafdavid23/vpc-cidr-manager/internal/aws"
	"github.com/asafdavid23/vpc-cidr-manager/internal/logging"
	"github.com/spf13/cobra"
)

// releaseCidrCmd represents the releaseCidr command
var releaseCidrCmd = &cobra.Command{
	Use:   "release-cidr",
	Short: "Release a CIDR block",
	Run: func(cmd *cobra.Command, args []string) {
		logLevel, err := cmd.Flags().GetString("log-level")
		cidr, err := cmd.Flags().GetString("cidr")

		logger := logging.NewLogger(logLevel)

		logger.Debug("Initializing DynamoDB client")
		client, err := internalAws.GetDynamoDBClient()

		if err != nil {
			logger.Fatal(err)
		}

		logger.Debug("Releasing CIDR block")
		err = internalAws.ReleaseCidr(client, cidr, logger)

		if err != nil {
			logger.Fatal(err)
		}

		logger.Infof("%s CIDR block released successfully", cidr)
	},
}

func init() {
	rootCmd.AddCommand(releaseCidrCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// releaseCidrCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// releaseCidrCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	releaseCidrCmd.Flags().StringP("cidr", "c", "", "The CIDR block to release")
	releaseCidrCmd.MarkFlagRequired("cidr")
	releaseCidrCmd.Flags().StringP("log-level", "l", "info", "The log level to use")
}
