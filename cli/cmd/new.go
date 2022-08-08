/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"fmt"
	"time"

	"github.com/gookit/color"
	"github.com/ostamand/aqualog/api"
	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Log a new param to Aqualog",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		paramType, _ := cmd.Flags().GetString("type")
		value, _ := cmd.Flags().GetFloat64("value")
		timeString, _ := cmd.Flags().GetString("timestamp")

		authChan := make(chan error)

		go func() {
			authChan <- aqualog.RenewTokenIf()
		}()

		req := api.CreateParamRequest{
			ParamType: paramType,
			Value:     value,
		}

		if timeString != "" {
			ts, err := time.ParseInLocation("2006-01-02 15:04:05", timeString, time.Local)
			if err != nil {
				fmt.Println(err)
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

		err := aqualog.CreateParam(req)
		if err != nil {
			if errors.Is(err, ErrNeedToLogin) {
				color.Error.Println("Please login to Aqualog using `aqualog login`.")
			} else {
				color.Error.Println("Param logging failed.")
			}
			return
		}
		color.Info.Println("Param logging successful.")
	},
}

func init() {
	paramCmd.AddCommand(newCmd)

	newCmd.Flags().StringP("type", "t", "", "type of param")
	newCmd.MarkFlagRequired("type")

	newCmd.Flags().Float64P("value", "v", 0, "the param value")
	newCmd.MarkFlagRequired("value")

	newCmd.Flags().String("timestamp", "", "when was the measurement taken")
}
