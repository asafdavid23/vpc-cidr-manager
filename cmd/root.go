/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var Version string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "vpc-cidr-manager",
	Short: "A CLI tool to manage VPC CIDR block reservations in DynamoDB",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		if versionFlag, _ := cmd.Flags().GetBool("version"); versionFlag {
			fmt.Println("vpc-cidr-manager version:", Version)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.vpc-cidr-manager.yaml)")
	rootCmd.PersistentFlags().String("log-level", "info", "Set the log level (debug, info, warn, error, fatal)")
	rootCmd.PersistentFlags().BoolP("version", "v", false, "Display the version of this CLI tool")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
