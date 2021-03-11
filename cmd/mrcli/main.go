package main

import (
	"github.com/danielpieper/mrcli-go/internal/commands"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{Use: "app"}
	rootCmd.AddCommand(commands.OverviewCmd)
	rootCmd.AddCommand(commands.ApproverCmd)
	rootCmd.AddCommand(commands.ProjectCmd)
	rootCmd.Execute()
}
