package main

import (
	"github.com/danielpieper/mrcli-go/cmd"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{Use: "app"}
	rootCmd.AddCommand(cmd.OverviewCmd)
	rootCmd.AddCommand(cmd.ApproverCmd)
	rootCmd.Execute()
}
