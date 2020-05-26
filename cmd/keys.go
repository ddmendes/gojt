package cmd

import (
	"fmt"
	"os"

	"github.com/ddmendes/gojt/jsondoc"
	"github.com/spf13/cobra"
)

var keysCmd = &cobra.Command{
	Use:     "keys",
	Aliases: []string{"k"},
	Short:   "Print keys under a given path",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		document, err := document.Get(path)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		keys, err := document.GetKeys()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		keysDoc := jsondoc.Wrap(keys)
		output, err := keysDoc.Marshal(true)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(string(output))
	},
}

func init() {
	rootCmd.AddCommand(keysCmd)
}
