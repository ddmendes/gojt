package cmd

import (
	"fmt"
	"os"

	"github.com/ddmendes/gojt/jsdoc"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gojt",
	Short: "GOJT is a utility tool for parsing and reading json documents from terminal",
}

var document jsdoc.JSDoc

func init() {
	if err := jsdoc.ReadPipedDoc(&document); err != nil {
		fmt.Println("ERROR: Failed to read JSON document.", err)
		panic(err)
	}
}

// Execute gojt CLI command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
