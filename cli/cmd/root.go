package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var aqualog aqualogAPI

const promptFailedMessage = "Prompt failed"

var rootCmd = &cobra.Command{
	Use:   "aqualog",
	Short: "Reef aquarium logging made easy, fun and free",
	Long:  ``,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	aqualog = NewAqualogAPI()
	aqualog.LoadAuth() // get auth from local folder if it exits

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
