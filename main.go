package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const (
	pattern     = "trw-files-*"
	zipFileName = "trw.zip"
	txtFileName = "trw.txt"
)

func main() {
	tempDirPath, err := os.MkdirTemp("", pattern)
	if err != nil {
		log.Fatal(err)
	}
	// deletes directory when main function exits
	// probably needs to change because im doing a schedular for the main function
	defer os.RemoveAll(tempDirPath)

	latestURL, err := GetLatestTaxRollURL()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("This week's target URL is:", latestURL)

	zipPath := filepath.Join(tempDirPath, zipFileName)
	err = DownloadFile(latestURL, zipPath)
	if err != nil {
		log.Println("Download failed:", err)
	}

	txtPath := filepath.Join(tempDirPath, txtFileName)
	err = UnzipFile(zipPath, txtPath)
	if err != nil {
		log.Println("Unzip failed:", err)
	}
}
