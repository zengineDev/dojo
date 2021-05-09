package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "dojo",
	Short: "The command line for the dojo framework",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		fmt.Println("TODO")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
