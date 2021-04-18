package main

import (
	"github.com/fatih/color"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "email-service version",
	Long:  `Version of the last release of email-service`,
	Run: func(cmd *cobra.Command, args []string) {
		c := color.New(color.FgBlack).Add(color.BgYellow).Add(color.Underline)
		c.Println("email-service v0.1")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
