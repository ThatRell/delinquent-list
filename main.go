package main

import (
	"fmt"
)

// downloadURL may need to change depending on if link changes
// in downloadfile, have a guaranteed way to get the zip file even if name changes
const (
	zipURL    = "https://www.dallascounty.org/Assets/uploads/docs/tax/trw/trwfile.725050.zip"
	zipPath   = "trw.zip"
	outputDir = "./data"
)

func main() {
	// os.MkdirAll(outputDir, os.ModePerm)
	// err := DownloadFile(zipURL, zipPath)
	// if err != nil {
	// 	fmt.Print("error downloading file")
	// 	return
	// }

	outputPath, err := ExtractFromZip(zipPath, outputDir)
	if err != nil {
		fmt.Print("error extracting zip")
		return
	}

	date, err := ExtractDate(outputPath)
	if err != nil {
		fmt.Print("error getting date")
		return
	}

	fmt.Println("Date: ", date)

	fmt.Println(outputPath)
	ParseFile(outputPath, outputDir, date)
}
