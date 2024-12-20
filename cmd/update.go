package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	getter "github.com/hashicorp/go-getter"
	"github.com/spf13/cobra"
)

var (
	updateCmd *cobra.Command
)

func runUpdate(version string) error {
	versionParam := "latest/download/dg"
	if version != "" {
		versionParam = fmt.Sprintf("download/%s/dg", version)
	}

	url := fmt.Sprintf("https://github.com/cognite-analytics/dbt-go/releases/%s", versionParam)
	fmt.Printf("Updating CLI tool from %s...\n", url)

	tempFile := filepath.Join(os.TempDir(), "dg")
	if runtime.GOOS == "windows" {
		tempFile += ".exe"
	}

	err := getter.GetFile(tempFile, url)
	if err != nil {
		return fmt.Errorf("failed to download file: %v", err)
	}

	targetDir := getInstallPath()
	targetPath := filepath.Join(targetDir, "dg")
	if runtime.GOOS == "windows" {
		targetPath += ".exe"
	}

	err = os.Rename(tempFile, targetPath)
	if err != nil {
		return fmt.Errorf("failed to move binary to target directory: %v", err)
	}

	if runtime.GOOS != "windows" {
		err = os.Chmod(targetPath, 0755)
		if err != nil {
			return fmt.Errorf("failed to change file permissions: %v", err)
		}
	}

	fmt.Println("Update successful!")
	return nil
}

func getInstallPath() string {
	if runtime.GOOS == "windows" {
		path, err := exec.LookPath(os.Args[0])
		if err == nil {
			return filepath.Dir(path)
		}
		return "."
	}
	return "/usr/local/bin"
}
