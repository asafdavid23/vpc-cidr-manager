/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/asafdavid23/vpc-cidr-manager/internal/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Version = "dev"
var cfgFile string

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
	initConfig()
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.vpc-cidr-manager.yaml)")
	rootCmd.PersistentFlags().String("log-level", "info", "Set the log level (debug, info, warn, error, fatal)")
	rootCmd.PersistentFlags().Bool("version", false, "Display the version of this CLI tool")
	rootCmd.PersistentFlags().String("output", "table", "Output type table/json/yaml")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	viper.BindPFlag("global.logLevel", rootCmd.PersistentFlags().Lookup("log-level"))
	viper.BindPFlag("global.output", rootCmd.PersistentFlags().Lookup("output"))
	viper.BindPFlag("global.region", rootCmd.PersistentFlags().Lookup("region"))
}

func initConfig() {
	logLevel := viper.GetString("global.logLevel")
	logger := logging.NewLogger(logLevel)

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// Use the default config file
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("./config/")
	}

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		// Ignore "file not found" errors; log other errors
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			logger.Fatalf("Error reading config file: %v", err)
		}
	}

	// Log the config file being used
	if viper.ConfigFileUsed() != "" {
		logger.Printf("Using config file: %s", viper.ConfigFileUsed())
	}
}
