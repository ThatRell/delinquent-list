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

func ExtractFromZip(zipPath, outputDir string) (string, error) {

	r, err := zip.OpenReader(zipPath)
	if err != nil {
		fmt.Println("open zip error")
		return "", err
	}
	defer r.Close()

	for _, f := range r.File {
		base := filepath.Base(f.Name)

		if !strings.HasPrefix(strings.ToLower(base), "flat404") {
			continue
		}

		rc, err := f.Open()
		if err != nil {
			fmt.Println("f open error")
			return "", err
		}
		defer rc.Close()

		outPath := filepath.Join(outputDir, base)

		outFile, err := os.Create(outPath)
		if err != nil {
			fmt.Println("outfile create error")
			return "", err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()

		if err != nil {
			fmt.Println("outfile copy error")
			return "", err
		}

		return outPath, nil
	}

	fmt.Println("no flat404 file found")
	return "", os.ErrNotExist
}
