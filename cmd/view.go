package cmd

import (
	lib "emr-idle-check/lib"
	"fmt"

	"github.com/spf13/cobra"
)

var ViewCMD = &cobra.Command{
	Use:   "view",
	Short: "view a snapshot of the current idle-check metrics",
	Run: func(cmd *cobra.Command, args []string) {
		idleMetrics := lib.GetIdleCheckMetrics()
		fmt.Println(idleMetrics)
	},
}
