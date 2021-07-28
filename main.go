package main

import (
	"fmt"
	"os"

	cmd "emr-idle-check/cmd"

	"github.com/spf13/cobra"
)

// Version of emr-idle-check. Overwritten during build
var version = "development"

var rootCmd = &cobra.Command{
	Use:   "emr-idle-check",
	Short: "Send EMR activity metrics to CloudWatch",
}

func main() {
	Execute(version)
}

// Import other commands
func init() {
	cmd.Import(rootCmd)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	rootCmd.Version = version
	rootCmd.SetVersionTemplate(`{{printf "%s" .Version}}
`)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
