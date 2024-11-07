package cmd

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
	"golang.org/x/mod/semver"

	"github.com/Shackelford-Arden/hctx/build"
	"github.com/Shackelford-Arden/hctx/internal/github"
)

func ShowPath(ctx *cli.Context) error {

	fullPath, _ := os.Executable()

	fmt.Println(fullPath)

	return nil
}

func SelfUpdate(ctx *cli.Context) error {

	gh := github.NewClient("", "")
	gt := os.Getenv("GITHUB_TOKEN")

	if gt != "" {
		gh.SetToken(gt)
	}

	latestVersion, lvError := gh.GetLatestRelease()
	if lvError != nil {
		return fmt.Errorf("Failed to get latest version: %s", lvError.Error())
	}

	if semver.Compare(build.Version, latestVersion.Version()) != -1 {
		return fmt.Errorf("current version is either already the latest, or we failed to parse versions correctly")
	}

	fmt.Println(fmt.Sprintf("Downloading %s...", latestVersion.Version()))

	// Identify the download URL
	downloadUrl := latestVersion.TarballUrl()

	if downloadUrl == "" {
		return fmt.Errorf("release information didn't include an appropriate file to download.")
	}

	// Download the release
	resp, err := http.Get(latestVersion.TarballUrl())
	if err != nil {
		return fmt.Errorf("failed to download the release: %w", err)
	}
	defer resp.Body.Close()

	// Create a temporary file to store the tarball
	downloadDir, err := os.MkdirTemp(os.TempDir(), "hctx")
	if err != nil {
		return fmt.Errorf("failed to create a temporary directory so download the latest version: %w", err)
	}
	defer os.RemoveAll(downloadDir)

	tmpFile, err := os.CreateTemp(downloadDir, "release-*.tar.gz")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	// Copy the downloaded content to the temporary file
	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save %s: %w", tmpFile.Name(), err)
	}

	// Rewind the file for reading
	_, err = tmpFile.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("failed to rewind temporary file: %w", err)
	}

	// Open and iterate through the tarball
	gr, err := gzip.NewReader(tmpFile)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gr.Close()

	tr := tar.NewReader(gr)

	// Get the current executable path so we know what needs to be replaced.
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get current executable path: %w", err)
	}

	// Find and extract the binary from within the tar.gz file.
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar: %w", err)
		}

		if header.Typeflag == tar.TypeReg && filepath.Base(header.Name) == "hctx" {
			// Create a temporary file for the new binary
			newBinary, err := os.CreateTemp(filepath.Dir(execPath), "new-binary-*")
			if err != nil {
				return fmt.Errorf("failed to create temporary file for new binary: %w", err)
			}
			defer os.Remove(newBinary.Name())

			// Copy the binary content
			_, err = io.Copy(newBinary, tr)
			if err != nil {
				return fmt.Errorf("failed to copy new binary: %w", err)
			}

			// Close the new binary file
			newBinary.Close()

			// Get the current executable's permissions
			info, err := os.Stat(execPath)
			if err != nil {
				return fmt.Errorf("failed to get current executable info: %w", err)
			}

			// Set the permissions on the new binary
			err = os.Chmod(newBinary.Name(), info.Mode())
			if err != nil {
				return fmt.Errorf("failed to set permissions on new binary: %w", err)
			}

			// Rename the new binary to replace the current executable
			err = os.Rename(newBinary.Name(), execPath)
			if err != nil {
				return fmt.Errorf("failed to replace current executable: %w", err)
			}

			return nil
		}
	}

	fmt.Println("Update successful!")

	return nil

}
