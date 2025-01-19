/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	internalAws "github.com/asafdavid23/vpc-cidr-manager/internal/aws"
	"github.com/asafdavid23/vpc-cidr-manager/internal/logging"
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

		logger.Debug("Initializing DynamoDB client")
		client, err := internalAws.GetDynamoDBClient()
		if err != nil {
			logger.Fatal(err)
		}

		err = internalAws.CreateDynamoDBTable(client, tableName)

		if err != nil {
			logger.Fatal(err)
		}

		logger.Info("DynamoDB table created successfully")
	},
}

func init() {
	rootCmd.AddCommand(createTableCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createTableCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createTableCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	createTableCmd.Flags().StringP("name", "n", "", "The name of the table to create")
}
