/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Aqualog service",
	Long:  `A successful login is required in order to execute any Aqualog CLI commands`,
	Run: func(cmd *cobra.Command, args []string) {
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		err := aqualog.Login(username, password)
		if err != nil {
			if errors.Is(err, ErrWrongPassword) {
				color.Error.Println("Wrong password. Try again using `aqualog login`.")
			} else if errors.Is(err, ErrUserNotFound) {
				color.Error.Println("User does not exists.") // TODO add info how to create a new user
			} else {
				color.Error.Println("Login failed. Try again.")
			}
			return
		}
		color.Info.Println("Login successful.")
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringP("password", "p", "", "User password")
	loginCmd.MarkFlagRequired("password")

	loginCmd.Flags().StringP("username", "u", "", "Username")
	loginCmd.MarkFlagRequired("username")
}
