package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
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

func IsOlder(dateA, dateB string) bool {
	a, _ := time.Parse("20060102", dateA)
	b, _ := time.Parse("20060102", dateB)

	return a.Before(b)
}

func ParseFile(outputPath string, outputDir string, currdate string) {

	csv_outputPath := filepath.Join(outputDir, "delinquent_properties.csv")

	csv_file, _ := os.Create(csv_outputPath)
	writer := csv.NewWriter(csv_file)
	defer writer.Flush()

	writer.Write([]string{
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
	})

	file, err := os.Open(outputPath)
	if err != nil {
		fmt.Println("couldn't open file to parse")
		panic(err)
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		rec := ParseLine(line)

		// filter based on delinquency
		// if DUE_DATE is before the YYYYMMDD in the output file name
		// then it's a delinquent account
		// the csv file is large so i think there's
		// something wrong with my filtering
		// also ints are displaying weird
		if IsOlder(rec.Due_Date, currdate) {
			writer.Write([]string{
				rec.Account,
				rec.Year,
				rec.Jurisdiction,
				rec.Tax_Unit_Acct,
				rec.Levy,
				rec.Homestead,
				rec.Over65,
				rec.Veteran,
				rec.Disabled,
				rec.AG,
				rec.Date_Paid,
				rec.Due_Date,
				rec.Omit_Flag,
				rec.Levy_Balance,
				rec.Suit,
				rec.Causeno,
				rec.Bankcode,
				rec.BankruptNo,
				rec.Attorney,
				rec.Court_Cost,
				rec.Abstract_Fee,
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
				rec.Parcel_No,
				rec.Parcel_Name,
				rec.Payment_Agreement,
				rec.Total_Amount_Due,
				rec.Total_Amount_Due_30,
				rec.Total_Amount_Due_60,
				rec.Total_Amount_Due_90,
				rec.Amount_Indicator,
			})
		}
	}
}

func ParseLine(line string) Record {
	// start_index is start_position - 1
	return Record{
		Account:             strings.TrimSpace(line[0:34]),
		Year:                strings.TrimSpace(line[34:38]),
		Jurisdiction:        strings.TrimSpace(line[38:42]),
		Tax_Unit_Acct:       strings.TrimSpace(line[42:76]),
		Levy:                strings.TrimSpace(line[76:87]),
		Homestead:           strings.TrimSpace(line[87:88]),
		Over65:              strings.TrimSpace(line[88:89]),
		Veteran:             strings.TrimSpace(line[89:90]),
		Disabled:            strings.TrimSpace(line[90:91]),
		AG:                  strings.TrimSpace(line[91:92]),
		Date_Paid:           strings.TrimSpace(line[92:100]),
		Due_Date:            strings.TrimSpace(line[100:108]),
		Omit_Flag:           strings.TrimSpace(line[108:110]),
		Levy_Balance:        strings.TrimSpace(line[110:121]),
		Suit:                strings.TrimSpace(line[121:122]),
		Causeno:             strings.TrimSpace(line[122:162]),
		Bankcode:            strings.TrimSpace(line[162:163]),
		BankruptNo:          strings.TrimSpace(line[163:203]),
		Attorney:            strings.TrimSpace(line[203:204]),
		Court_Cost:          strings.TrimSpace(line[204:211]),
		Abstract_Fee:        strings.TrimSpace(line[211:218]),
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
		Zip:                 strings.TrimSpace(line[427:439]),
		Roll_Code:           strings.TrimSpace(line[439:440]),
		Parcel_No:           strings.TrimSpace(line[440:448]),
		Parcel_Name:         strings.TrimSpace(line[448:488]),
		Payment_Agreement:   strings.TrimSpace(line[488:489]),
		Total_Amount_Due:    strings.TrimSpace(line[489:500]),
		Total_Amount_Due_30: strings.TrimSpace(line[500:511]),
		Total_Amount_Due_60: strings.TrimSpace(line[511:522]),
		Total_Amount_Due_90: strings.TrimSpace(line[522:533]),
		Amount_Indicator:    strings.TrimSpace(line[533:534]),
	}
}
