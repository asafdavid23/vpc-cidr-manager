/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// generateTemplateCmd represents the generateTemplate command
var generateTemplateCmd = &cobra.Command{
	Use:   "generate-template",
	Short: "Generate IaaC templates",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("generateTemplate called")
	},
}

func init() {
	// rootCmd.AddCommand(generateTemplateCmd)
	iaacCmd.AddCommand(generateTemplateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateTemplateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateTemplateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
