package main

import (
	"fmt"
	"os"
)

// downloadURL may need to change depending on if link changes
// in downloadfile, have a guaranteed way to get the zip file even if name changes
const (
	zipURL    = "https://www.dallascounty.org/Assets/uploads/docs/tax/trw/trwfile.725050.zip"
	zipPath   = "trw.zip"
	outputDir = "./data"
)

func main() {
	os.MkdirAll(outputDir, os.ModePerm)
	err := DownloadFile(zipURL, zipPath)
	if err != nil {
		fmt.Print("error downloading file")
		return
	}

	outputPath, err := ExtractFromZip(zipPath, outputDir)
	fmt.Println("downloaded and extracted")
	fmt.Println(outputPath)
}
