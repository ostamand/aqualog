/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gookit/color"
	"github.com/manifoldco/promptui"
	"github.com/ostamand/aqualog/api"
	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Log a new param to Aqualog",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		paramType, _ := cmd.Flags().GetString("type")
		value, _ := cmd.Flags().GetFloat64("value")
		timeString, _ := cmd.Flags().GetString("timestamp")

		// parallel call to renew token
		authChan := make(chan error)
		go func() {
			authChan <- aqualog.RenewTokenIf()
		}()

		// check if param type was provided
		if paramType == "" {
			prompt := promptui.Select{
				Label: "Select type",
				Items: []string{"Phosphate", "Nitrate", "PH", "Alkalinity", "Calcium", "Magnesium"},
			}
			_, paramType, err = prompt.Run()
			paramType = strings.ToLower(paramType) // use all lower cases for API call
			if err != nil {
				color.Error.Println(promptFailedMessage)
				return
			}
		}

		// check if value was provided
		if value == 0 {
			prompt := promptui.Prompt{
				Label: "Number",
				Validate: func(s string) error {
					_, err := strconv.ParseFloat(s, 64)
					return err
				},
			}
			result, err := prompt.Run()
			if err != nil {
				color.Error.Println(promptFailedMessage)
				return
			}
			value, _ = strconv.ParseFloat(result, 64)
		}

		req := api.CreateParamRequest{
			ParamType: paramType,
			Value:     value,
		}

		// check if timestamp was provided
		if timeString != "" {
			ts, err := time.ParseInLocation("2006-01-02 15:04:05", timeString, time.Local)
			if err != nil {
				color.Error.Println("Wrong timestamp format, use: YYYY-MM-DD HH:MM:SS")
				return
			}
			req.Timestamp = ts
		}

		// check access token first
		authErr := <-authChan
		if authErr != nil {
			fmt.Println(authErr)
			color.Error.Println("Please login to Aqualog using `aqualog login`.")
			return
		}

		err = aqualog.CreateParam(req)
		if err != nil {
			if errors.Is(err, ErrNeedToLogin) {
				color.Error.Println("Please login to Aqualog using `aqualog login`.")
			} else {
				color.Error.Println("Param logging failed.")
			}
			return
		}
		color.Info.Println("Param logged successfully")
	},
}

func init() {
	paramCmd.AddCommand(newCmd)

	newCmd.Flags().StringP("type", "t", "", "type of param")
	newCmd.Flags().Float64P("value", "v", 0, "the param value")
	newCmd.Flags().String("timestamp", "", "when was the measurement taken")
}
