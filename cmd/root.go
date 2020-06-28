package cmd

import (
	"fmt"
	"os"

	"github.com/ddmendes/gojt/jsondoc"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gojt",
	Short: "GOJT is a utility tool for parsing and reading json documents from terminal",
}

var document jsondoc.JSONDoc

func init() {
	var err error
	jsonReader := jsondoc.NewJSONReader(os.Stdin)
	document, err = jsonReader.ReadJSON()
	if err != nil {
		fmt.Println("ERROR: Failed to read JSON document.", err)
		printErrorAndQuit(err, 1)
	}
}

func printErrorAndQuit(err error, exitCode int) {
	fmt.Println(err)
	os.Exit(exitCode)
}

// Execute gojt CLI command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
