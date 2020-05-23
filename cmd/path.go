package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var pathCmd = &cobra.Command{
	Use:     "path",
	Aliases: []string{"p"},
	Short:   "Print object on a given path",
	Long:    "Print object on a given path",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("keysCmd")
		fmt.Println("cmd:", cmd)
		fmt.Println("args:", args)
	},
}

func init() {
	rootCmd.AddCommand(pathCmd)
}
