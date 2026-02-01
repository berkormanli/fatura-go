package internal

import "regexp"

var patterns = map[string]*regexp.Regexp{
	"invoiceNumber": regexp.MustCompile(`^(GIB)[a-zA-Z0-9]{13}$`),
	"gtipCode":      regexp.MustCompile(`^[0-9]{12}$`),
	"date":          regexp.MustCompile(`^(0[1-9]|1[0-9]|2[0-9]|3(0|1))/(0[1-9]|1[0-2])/\d{4}$`),
	"time":          regexp.MustCompile(`^(?:2[0-3]|[01][0-9]):[0-5][0-9]:[0-5][0-9]$`),
}

func ValidateInvoiceNumber(value string) bool {
	return patterns["invoiceNumber"].MatchString(value)
}

func ValidateGtipCode(value string) bool {
	return patterns["gtipCode"].MatchString(value)
}

func ValidateDate(value string) bool {
	return patterns["date"].MatchString(value)
}

func ValidateTime(value string) bool {
	return patterns["time"].MatchString(value)
}
