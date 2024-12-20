package cmd

import (
	"dg/style"
	"fmt"
	"os"
	"runtime"

	getter "github.com/hashicorp/go-getter"
	"github.com/spf13/cobra"
)

// Define the tool version
var version = "v0.1.0"

var rootCmd = &cobra.Command{
	Use:   "dg",
	Short: style.Dg.Render("\ndg is a cli written in GO to help improve the DX of Analytics Engineers using dbt"),
	Long:  style.Dg.Render("dg is a cli written in GO to help improve the DX of Analytics Engineers using dbt"),
	RunE:  executeRoot,
}

func init() {
	rootCmd.PersistentFlags().BoolP("version", "v", false, "Print the version number of the CLI tool")
	rootCmd.Flags().BoolP("update", "u", false, "Update the CLI tool to the latest version or a specified version")
}

func executeRoot(cmd *cobra.Command, args []string) error {
	if cmd.Flag("version").Changed {
		fmt.Printf("Current Version: %s\n", version)
		return nil
	}

	if cmd.Flag("update").Changed {
		var updateVersion string
		if len(args) > 0 {
			updateVersion = args[0]
		}
		return runUpdate(updateVersion)
	}

	cmd.Help()
	return nil
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runUpdate(version string) error {
	versionParam := "latest/download/dg"

	if version != "" {
		versionParam = fmt.Sprintf("download/%s/dg", version)
	}

	url := fmt.Sprintf("https://github.com/cognite-analytics/dbt-go/releases/%s", versionParam)

	fmt.Printf("Updating CLI tool from %s...\n", url)

	tempFile := "/tmp/dg"
	if runtime.GOOS == "windows" {
		tempFile = "dg.exe"
	}

	err := getter.GetFile(tempFile, url)
	if err != nil {
		return fmt.Errorf("failed to download file: %v", err)
	}

	err = os.Rename(tempFile, os.Args[0])
	if err != nil {
		return fmt.Errorf("failed to replace the binary: %v", err)
	}

	if runtime.GOOS != "windows" {
		err = os.Chmod(os.Args[0], 0755)
		if err != nil {
			return fmt.Errorf("failed to change the file permissions: %v", err)
		}
	}

	fmt.Println("Update successful!")
	return nil
}
