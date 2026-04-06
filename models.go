package main

type Record struct {
	Account             string
	Year                int
	Jurisdiction        int
	Tax_Unit_Acct       string
	Levy                int
	Homestead           string
	Over65              string
	Veteran             string
	Disabled            string
	AG                  string
	Date_Paid           string
	Due_Date            string
	Omit_Flag           string
	Levy_Balance        int
	Suit                string
	Causeno             string
	Bankcode            string
	BankruptNo          string
	Attorney            string
	Court_Cost          int
	Abstract_Fee        int
	Deferral            string
	Billsupp            string
	Split_PMTFlag       string
	Category_Code       string
	Owner               string
	Address2            string
	Address3            string
	Address4            string
	City                string
	State               string
	Zip                 string
	Roll_Code           string
	Parcel_No           int
	Parcel_Name         string
	Payment_Agreement   string
	Total_Amount_Due    int
	Total_Amount_Due_30 int
	Total_Amount_Due_60 int
	Total_Amount_Due_90 int
	Amount_Indicator    string
}

type Entry struct {
	Account          string
	Owner            string
	Address2         string
	Address3         string
	Address4         string
	City             string
	State            string
	Zip              string
	Total_Amount_Due int
}

func (rec *Record) ToEntry() Entry {
	return Entry{
		Account:          rec.Account,
		Owner:            rec.Owner,
		Address2:         rec.Address2,
		Address3:         rec.Address3,
		Address4:         rec.Address4,
		City:             rec.City,
		State:            rec.State,
		Zip:              rec.Zip,
		Total_Amount_Due: rec.Total_Amount_Due,
	}
}
