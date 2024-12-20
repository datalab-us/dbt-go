package cmd

import (
	"dg/style"
	"dg/version"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dg",
	Short: style.Dg.Render("\ndg is a cli written in GO to help improve the DX of Analytics Engineers using dbt"),
	Long:  style.Dg.Render("dg is a cli written in GO to help improve the DX of Analytics Engineers using dbt"),
	RunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Flag("version").Changed {
			checkVersion()
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

func checkVersion() {
	currentVersion := version.Version
	latestVersion, err := getLatestVersion()
	if err != nil {
		fmt.Printf("Error checking latest version: %v\n", err)
		return
	}

	if currentVersion < latestVersion {
		fmt.Printf("Your version (%s) is out of date. Please run 'dg upgrade' to update to the latest version (%s).\n", currentVersion, latestVersion)
	} else {
		fmt.Println(style.LightGray.Render(fmt.Sprintf("\nCurrent Version: %s", currentVersion)))
		fmt.Println(style.Dg.Render(fmt.Sprintf("\nYou are up to date!")))
	}
}

func getLatestVersion() (string, error) {
	resp, err := http.Get("https://api.github.com/repos/cognite-analytics/dbt-go/releases/latest")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	return release.TagName, nil
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
