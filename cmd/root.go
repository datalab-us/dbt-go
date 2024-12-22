package cmd

import (
	"dg/style"
	"dg/version"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings" // Add this import

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
		if cmd.Flag("info").Changed {
			runInfo(cmd, args)
			return nil
		}
		cmd.Help()
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().BoolP("version", "v", false, "Print the version number of the CLI tool")
	rootCmd.PersistentFlags().BoolP("info", "i", false, "Show Additional Developer Information About dbt-go")

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
		fmt.Printf("Your version (%s) is out of date. Would you like to upgrade to the latest version (%s)? (yes/no): ", currentVersion, latestVersion)
		var response string
		fmt.Scanln(&response)
		if response == "yes" {
			runUpdate(latestVersion)
		} else {
			fmt.Println(style.LightGray.Render(fmt.Sprintf("\nCurrent Version: %s", currentVersion)))
			fmt.Println(style.Red.Render("\nYou are out of date!"))
		}
	} else {
		fmt.Println(style.LightGray.Render(fmt.Sprintf("\nCurrent Version: %s", currentVersion)))
		fmt.Println(style.Dg.Render("\nYou are up to date!"))
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

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func stripANSI(str string) string {
	return ansiRegex.ReplaceAllString(str, "")
}

func centerText(text string, width int) string {
	stripped := stripANSI(text)
	if len(stripped) >= width {
		return text
	}
	spaces := width - len(stripped)
	left := spaces / 2
	return strings.Repeat(" ", left) + text
}

func runInfo(cmd *cobra.Command, args []string) error {
	asciiArt := `
--------------------------
--------------------------
--------*@@*--------------
------===+@#----===-------
---=#@@@@%@#--+%@@@%%@%=--
--=@@*==+%@#-+@%===#@@=---
--+@*:--:+@*:%@=:-:-@%----
--=%@*++*@@%++@@*++%@%----
---=*%@@%#%%#-=#%%%#@@----
----------------=+=*@#----
---------------#@@@%*-----
--------------------------
--------------------------
`
	copyright := "Copyright © 2024"
	contact := "Matthew Skinner -- matthew@skinnerdev.com"
	github := "github.com/cognite-analytics/dbt-go"
	width := 80

	centeredCopyright := centerText(copyright, width)
	centeredContact := centerText(contact, width)
	centeredLink := centerText(github, width)
	lines := strings.Split(asciiArt, "\n")
	centeredLines := make([]string, len(lines))
	for i, line := range lines {
		centeredLines[i] = centerText(line, width)
	}

	centeredAsciiArt := strings.Join(centeredLines, "\n")
	styledAsciiArt := style.Dg.Render(centeredAsciiArt)

	fmt.Printf(`%s

%s
%s
%s
`, styledAsciiArt, centeredCopyright, centeredContact, centeredLink)
	return nil
}
