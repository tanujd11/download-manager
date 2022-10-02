/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "downloader",
	Short: "A concurrent downloader in golang",
	Long: `downloader is a CLI tool to concurrently handle download
	of files of huge size concurrently to make it fast`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("download manager called")
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

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.download-manager.yaml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
