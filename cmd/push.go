package cmd

import (
	lib "emr-idle-check/lib"
	"fmt"

	"github.com/spf13/cobra"
)

var PushCMD = &cobra.Command{
	Use:   "push",
	Short: "push the idle-check metrics to AWS CloudWatch",
	Run: func(cmd *cobra.Command, args []string) {
		err := lib.PushIdleCheckMetricsToCloudWatch()
		if err != nil {
			fmt.Println(err)
		}
	},
}
