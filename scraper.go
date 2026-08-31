package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	fileTitle       = "TRW File"
	sampleFileTitle = "TRW Sample File"
)

func GetLatestTaxRollURL() (string, error) {
	baseURL := "https://www.dallascounty.org"
	landingPageURL := baseURL + "/departments/tax/tax-roll.php"

	res, err := http.Get(landingPageURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch landing page: %w", err)
	}
	defer res.Body.Close()

	// 200 OK is successful response status code
	if res.StatusCode != 200 {
		return "", fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %w", err)
	}

	var fileURL string

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		title, tExists := s.Attr("title")

		// samplefiletitle for testing
		if exists && tExists && title == sampleFileTitle && strings.HasSuffix(href, ".zip") {
			fileURL = href
		}
	})

	if fileURL == "" {
		return "", fmt.Errorf("could not find the download link on the page")
	}

	// for cases where URL is relative (when href doesn't contain base URL)
	parsedURL, err := url.Parse(fileURL)
	if err != nil {
		return "", err
	}

	if !parsedURL.IsAbs() {
		fileURL = baseURL + fileURL
	}

	return fileURL, nil
}
