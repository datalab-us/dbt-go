package cmd

import (
	"dg/style"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dg",
	Short: style.Dg.Render("\ndg is a cli written in GO to help improve the DX of Analytics Engineers using dbt"),
	Long:  style.Dg.Render("dg is a cli written in GO to help improve the DX of Analytics Engineers using dbt"),
	RunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Flag("version").Changed {
			fmt.Printf("Current Version: %s\n", version)
			return nil
		}
		cmd.Help()
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().BoolP("version", "v", false, "Print the version number of the CLI tool")

	updateCmd = &cobra.Command{
		Use:   "update [version]",
		Short: "Update the CLI tool to the latest or specified version",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var updateVersion string
			if len(args) > 0 {
				updateVersion = args[0]
			}
			return runUpdate(updateVersion)
		},
	}

	rootCmd.AddCommand(updateCmd)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
