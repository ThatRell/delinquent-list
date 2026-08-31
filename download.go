package main

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func DownloadFile(url, filepath string) error {
	fmt.Println("Downloading file")

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to request file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create local file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed during streaming: %w", err)
	}

	return nil
}

func UnzipFile(zipPath, destPath string) error {
	fmt.Println("Extracting zip")

	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %w", err)
	}
	defer r.Close()

	for _, f := range r.File {
		base := filepath.Base(f.Name)
		if !strings.HasPrefix(strings.ToLower(base), "flat404") {
			continue
		}

		rc, err := f.Open()
		if err != nil {
			return fmt.Errorf("f open error: %w", err)
		}
		defer rc.Close()

		out, err := os.Create(destPath)
		if err != nil {
			return fmt.Errorf("failed to create local file: %w", err)
		}
		defer out.Close()

		_, err = io.Copy(out, rc)
		if err != nil {
			return fmt.Errorf("failed during streaming: %w", err)
		}

		return nil
	}

	fmt.Println("no flat404 file found")
	return os.ErrNotExist
}
