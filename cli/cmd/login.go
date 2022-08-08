/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"fmt"

	"github.com/gookit/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func promptIfNeeded(value *string, label string, validate func(string) error) (err error) {
	if *value == "" {
		prompt := promptui.Prompt{
			Label:    label,
			Validate: validate,
		}
		*value, err = prompt.Run()
		if err != nil {
			return err
		}
	}
	return err
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Aqualog service",
	Long:  `A successful login is required in order to execute any Aqualog CLI commands`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")

		// username
		err = promptIfNeeded(&username, "Username", func(s string) error {
			if len(s) == 0 {
				return fmt.Errorf("must enter a username")
			}
			return nil
		})
		if err != nil {
			color.Error.Println(promptFailedMessage)
			return
		}

		// password
		err = promptIfNeeded(&password, "Password", func(s string) error {
			if len(s) < 6 {
				return fmt.Errorf("password must be more than 5 characters")
			}
			return nil
		})
		if err != nil {
			color.Error.Println(promptFailedMessage)
			return
		}

		err = aqualog.Login(username, password)

		if err != nil {
			if errors.Is(err, ErrWrongPassword) {
				color.Error.Println("Wrong password. Try again using `aqualog login`.")
			} else if errors.Is(err, ErrUserNotFound) {
				color.Error.Println("User does not exists.") // TODO add info how to create a new user
			} else {
				color.Error.Println("Login failed.")
			}
			return
		}
		color.Info.Println("Login successful.")
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().StringP("password", "p", "", "User password")
	loginCmd.Flags().StringP("username", "u", "", "Username")
}
