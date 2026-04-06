package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func ExtractDate(filepath string) (string, error) {
	filename := filepath[strings.Index(filepath, "\\")+1:]
	parts := strings.Split(filename, ".")

	if len(parts) < 3 {
		return "", fmt.Errorf("invalid filename format")
	}

	return parts[2], nil
}

func isBefore(dateA, dateB string) bool {
	a, _ := time.Parse("20060102", dateA)
	b, _ := time.Parse("20060102", dateB)
	return a.Before(b)
}

func OwesMoney(rec *Record) bool {
	return rec.Total_Amount_Due > 0
}

func IsNonResidential(name string) bool {

	words := strings.Split(strings.ToUpper(name), " ")
	has := map[string]bool{}
	for _, k := range words {
		has[k] = true
	}

	keywords := []string{
		"LLC", "INC", "BANK", "CITY", "COUNTY", "DISTRICT",
		"ISD", "UTILITY", "UTILITIES", "ROW", "PIPELINE",
	}

	for _, k := range keywords {
		if has[k] {
			return true
		}
	}

	return false
}

func filter(rec *Record) bool {
	if rec.State != "TX" {
		return false
	}
	if IsNonResidential(rec.Owner) {
		return false
	}
	if IsNonResidential(rec.Address2) {
		return false
	}
	if IsNonResidential(rec.Address3) {
		return false
	}
	if IsNonResidential(rec.Address4) {
		return false
	}

	return true
}

func ParseFile(outputPath string, outputDir string, currdate string) {

	csv_outputPath := filepath.Join(outputDir, "delinquent_properties.csv")

	csv_file, _ := os.Create(csv_outputPath)
	writer := csv.NewWriter(csv_file)
	defer writer.Flush()

	//writer.Write(RecordHeaders())
	writer.Write(EntryHeaders())

	file, err := os.Open(outputPath)
	if err != nil {
		fmt.Println("couldn't open file to parse")
		panic(err)
	}

	count := 0
	deliCount := 0

	seen := map[string]bool{}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		rec := ParseLine(line)
		count += 1

		// filter based on delinquency + other methods
		if !seen[rec.Account+rec.Owner+rec.Address2+rec.Address3+rec.Address4] && isBefore(rec.Due_Date, currdate) && OwesMoney(&rec) && filter(&rec) {

			deliCount += 1
			seen[rec.Account+rec.Owner+rec.Address2+rec.Address3+rec.Address4] = true
			entry := rec.ToEntry()

			//writer.Write(rec.RecordStrings())
			writer.Write(entry.EntryStrings())
		}
	}
	//fmt.Println(scanner.Err())
	fmt.Println(count, " records parsed")
	fmt.Println(deliCount, " delinquent records written")
}

func NumericField(line string, start int, end int) int {
	str := strings.TrimSpace(line[start:end])
	if str == "" {
		return 0
	}
	field, err := strconv.Atoi(str)
	if err != nil {
		fmt.Println("couldn't convert field to int")
		fmt.Println("error string: ", line[start:end])
		panic(err)
	}

	return field
}

func ParseLine(line string) Record {

	// zip issue fix
	zip := strings.TrimSpace(line[427:439])
	if len(zip) == 9 && zip[5:9] == "0000" {
		zip = zip[:5]
	}

	// start_index is start_position - 1
	return Record{
		Account:             strings.TrimSpace(line[0:34]),
		Year:                NumericField(line, 34, 38),
		Jurisdiction:        NumericField(line, 38, 42),
		Tax_Unit_Acct:       strings.TrimSpace(line[42:76]),
		Levy:                NumericField(line, 76, 87),
		Homestead:           strings.TrimSpace(line[87:88]),
		Over65:              strings.TrimSpace(line[88:89]),
		Veteran:             strings.TrimSpace(line[89:90]),
		Disabled:            strings.TrimSpace(line[90:91]),
		AG:                  strings.TrimSpace(line[91:92]),
		Date_Paid:           strings.TrimSpace(line[92:100]),
		Due_Date:            strings.TrimSpace(line[100:108]),
		Omit_Flag:           strings.TrimSpace(line[108:110]),
		Levy_Balance:        NumericField(line, 110, 121),
		Suit:                strings.TrimSpace(line[121:122]),
		Causeno:             strings.TrimSpace(line[122:162]),
		Bankcode:            strings.TrimSpace(line[162:163]),
		BankruptNo:          strings.TrimSpace(line[163:203]),
		Attorney:            strings.TrimSpace(line[203:204]),
		Court_Cost:          NumericField(line, 204, 211),
		Abstract_Fee:        NumericField(line, 211, 218),
		Deferral:            strings.TrimSpace(line[218:219]),
		Billsupp:            strings.TrimSpace(line[219:220]),
		Split_PMTFlag:       strings.TrimSpace(line[220:221]),
		Category_Code:       strings.TrimSpace(line[221:225]),
		Owner:               strings.TrimSpace(line[225:265]),
		Address2:            strings.TrimSpace(line[265:305]),
		Address3:            strings.TrimSpace(line[305:345]),
		Address4:            strings.TrimSpace(line[345:385]),
		City:                strings.TrimSpace(line[385:425]),
		State:               strings.TrimSpace(line[425:427]),
		Zip:                 zip,
		Roll_Code:           strings.TrimSpace(line[439:440]),
		Parcel_No:           NumericField(line, 440, 448),
		Parcel_Name:         strings.TrimSpace(line[448:488]),
		Payment_Agreement:   strings.TrimSpace(line[488:489]),
		Total_Amount_Due:    NumericField(line, 489, 500),
		Total_Amount_Due_30: NumericField(line, 500, 511),
		Total_Amount_Due_60: NumericField(line, 511, 522),
		Total_Amount_Due_90: NumericField(line, 522, 533),
		Amount_Indicator:    strings.TrimSpace(line[533:534]),
	}
}

func (rec *Record) RecordStrings() []string {
	return []string{
		rec.Account,
		strconv.Itoa(rec.Year),
		strconv.Itoa(rec.Jurisdiction),
		rec.Tax_Unit_Acct,
		strconv.Itoa(rec.Levy),
		rec.Homestead,
		rec.Over65,
		rec.Veteran,
		rec.Disabled,
		rec.AG,
		rec.Date_Paid,
		rec.Due_Date,
		rec.Omit_Flag,
		strconv.Itoa(rec.Levy_Balance),
		rec.Suit,
		rec.Causeno,
		rec.Bankcode,
		rec.BankruptNo,
		rec.Attorney,
		strconv.Itoa(rec.Court_Cost),
		strconv.Itoa(rec.Abstract_Fee),
		rec.Deferral,
		rec.Billsupp,
		rec.Split_PMTFlag,
		rec.Category_Code,
		rec.Owner,
		rec.Address2,
		rec.Address3,
		rec.Address4,
		rec.City,
		rec.State,
		rec.Zip,
		rec.Roll_Code,
		strconv.Itoa(rec.Parcel_No),
		rec.Parcel_Name,
		rec.Payment_Agreement,
		strconv.Itoa(rec.Total_Amount_Due),
		strconv.Itoa(rec.Total_Amount_Due_30),
		strconv.Itoa(rec.Total_Amount_Due_60),
		strconv.Itoa(rec.Total_Amount_Due_90),
		rec.Amount_Indicator,
	}
}

func RecordHeaders() []string {
	return []string{
		"account",
		"year",
		"jurisdiction",
		"tax-unit-acct",
		"levy",
		"homestead",
		"over65",
		"veteran",
		"disabled",
		"ag",
		"date-paid",
		"due-date",
		"omit-flag",
		"levy-balance",
		"suit",
		"causeno",
		"bankcode",
		"bankruptno",
		"attorney",
		"court-cost",
		"abstract-fee",
		"deferral",
		"billsupp",
		"split-pmtflag",
		"category-code",
		"owner",
		"address2",
		"address3",
		"address4",
		"city",
		"state",
		"zip",
		"roll-code",
		"parcel no.",
		"parcel name",
		"payment agreement",
		"tot_amt_due",
		"tot_amt_due-30",
		"tot_amt_due-60",
		"tot_amt_due-90",
		"amount indicator",
	}
}

func (e *Entry) EntryStrings() []string {
	return []string{
		e.Account,
		e.Owner,
		e.Address2,
		e.Address3,
		e.Address4,
		e.City,
		e.State,
		e.Zip,
		strconv.Itoa(e.Total_Amount_Due),
	}
}

func EntryHeaders() []string {
	return []string{
		"account",
		"owner",
		"address2",
		"address3",
		"address4",
		"city",
		"state",
		"zip",
		"amount due",
	}
}
