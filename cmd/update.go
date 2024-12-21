package cmd

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

var (
	updateCmd *cobra.Command
)

func runUpdate(version string) error {
	versionParam := "latest"
	if version != "" {
		versionParam = version
	}

	ext := "tar.gz"
	if runtime.GOOS == "windows" {
		ext = "zip"
	}

	url := fmt.Sprintf("https://github.com/cognite-analytics/dbt-go/releases/download/%s/dg_%s_%s_%s.%s", versionParam, versionParam, runtime.GOOS, runtime.GOARCH, ext)
	checksumURL := fmt.Sprintf("https://github.com/cognite-analytics/dbt-go/releases/download/%s/checksums.txt", versionParam)
	fmt.Printf("Updating CLI tool from %s...\n", url)

	tempFile := filepath.Join(os.TempDir(), fmt.Sprintf("dg.%s", ext))
	checksumFile := filepath.Join(os.TempDir(), "checksums.txt")

	err := downloadFile(tempFile, url)
	if err != nil {
		return fmt.Errorf("failed to download file: %v", err)
	}

	err = downloadFile(checksumFile, checksumURL)
	if err != nil {
		return fmt.Errorf("failed to download checksum file: %v", err)
	}

	err = validateChecksum(tempFile, checksumFile)
	if err != nil {
		return fmt.Errorf("checksum validation failed: %v", err)
	}

	targetDir := getInstallPath()
	if ext == "zip" {
		err = unzip(tempFile, targetDir)
	} else {
		err = untar(tempFile, targetDir)
	}
	if err != nil {
		return fmt.Errorf("failed to extract file: %v", err)
	}

	targetPath := filepath.Join(targetDir, "dg")
	if runtime.GOOS == "windows" {
		targetPath += ".exe"
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

func downloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func validateChecksum(file, checksumFile string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}
	computedChecksum := fmt.Sprintf("%x", h.Sum(nil))

	checksums, err := os.ReadFile(checksumFile)
	if err != nil {
		return err
	}

	for _, line := range strings.Split(string(checksums), "\n") {
		if strings.Contains(line, filepath.Base(file)) {
			expectedChecksum := strings.Fields(line)[0]
			if computedChecksum != expectedChecksum {
				return fmt.Errorf("checksum mismatch: expected %s, got %s", expectedChecksum, computedChecksum)
			}
			return nil
		}
	}
	return fmt.Errorf("checksum not found for file: %s", filepath.Base(file))
}

func unzip(src string, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, rc)

		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}

func untar(src string, dest string) error {
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()

	gr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gr.Close()

	tr := tar.NewReader(gr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		target := filepath.Join(dest, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, os.ModePerm); err != nil {
				return err
			}
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}
			f.Close()
		default:
			return fmt.Errorf("unknown type: %v in %s", header.Typeflag, header.Name)
		}
	}
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
