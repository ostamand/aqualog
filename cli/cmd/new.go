/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"

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
		//ts, _ := cmd.Flags().GetString("timestamp")

		err := aqualog.CreateParam(api.CreateParamRequest{
			ParamType: paramType,
			Value:     value,
		})
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

	//newCmd.Flags().StringP("timestamp", "ts", "", "when was the measurement taken")
}
