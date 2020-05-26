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
		path := args[0]
		document, err := document.Get(path)
		if err != nil {
			printErrorAndQuit(err, 1)
		}

		output, err := document.Marshal(true)
		if err != nil {
			printErrorAndQuit(err, 1)
		}
		fmt.Println(string(output))
	},
}

func init() {
	rootCmd.AddCommand(pathCmd)
}
