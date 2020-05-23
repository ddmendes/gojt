package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var keysCmd = &cobra.Command{
	Use:     "keys",
	Aliases: []string{"k"},
	Short:   "Print keys under a given path",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, k := range document.GetKeys() {
			fmt.Println(k)
		}
	},
}

func init() {
	rootCmd.AddCommand(keysCmd)
}
