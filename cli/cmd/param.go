/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// paramCmd represents the param command
var paramCmd = &cobra.Command{
	Use:   "param",
	Short: "Save a new parameter",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("param called")
	},
}

func init() {
	rootCmd.AddCommand(paramCmd)
}
