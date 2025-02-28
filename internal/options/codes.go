package options

import "sort"

const countryCodeReference = "https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2"

var validCountries = map[string]struct{}{
	"AD": {}, "AE": {}, "AF": {}, "AG": {}, "AI": {}, "AL": {}, "AM": {}, "AO": {}, "AR": {}, "AT": {},
	"AU": {}, "AX": {}, "AZ": {}, "BA": {}, "BD": {}, "BE": {}, "BG": {}, "BH": {}, "BI": {}, "BM": {},
	"BO": {}, "BQ": {}, "BR": {}, "BY": {}, "BZ": {}, "CA": {}, "CH": {}, "CL": {}, "CN": {}, "CO": {},
	"CW": {}, "CY": {}, "CZ": {}, "DE": {}, "DJ": {}, "DK": {}, "DM": {}, "DO": {}, "DZ": {}, "EC": {},
	"EE": {}, "EG": {}, "ES": {}, "EU": {}, "FI": {}, "FK": {}, "FO": {}, "FR": {}, "GB": {}, "GE": {},
	"GG": {}, "GI": {}, "GL": {}, "GM": {}, "GP": {}, "GR": {}, "GT": {}, "HK": {}, "HR": {}, "HU": {},
	"ID": {}, "IE": {}, "IL": {}, "IM": {}, "IN": {}, "IQ": {}, "IR": {}, "IS": {}, "IT": {}, "JE": {},
	"JO": {}, "JP": {}, "KE": {}, "KG": {}, "KH": {}, "KR": {}, "KW": {}, "KY": {}, "KZ": {}, "LB": {},
	"LI": {}, "LK": {}, "LR": {}, "LT": {}, "LU": {}, "LV": {}, "LY": {}, "MA": {}, "MC": {}, "MD": {},
	"ME": {}, "MH": {}, "MK": {}, "MO": {}, "MQ": {}, "MT": {}, "MU": {}, "MV": {}, "MX": {}, "MY": {},
	"MZ": {}, "NG": {}, "NI": {}, "NL": {}, "NO": {}, "NP": {}, "NZ": {}, "OM": {}, "PA": {}, "PE": {},
	"PH": {}, "PK": {}, "PL": {}, "PS": {}, "PT": {}, "PW": {}, "QA": {}, "RE": {}, "RO": {}, "RS": {},
	"RU": {}, "RW": {}, "SA": {}, "SC": {}, "SE": {}, "SG": {}, "SI": {}, "SK": {}, "SL": {}, "SM": {},
	"SN": {}, "SY": {}, "TG": {}, "TH": {}, "TJ": {}, "TK": {}, "TM": {}, "TN": {}, "TR": {}, "TW": {},
	"UA": {}, "UG": {}, "US": {}, "UZ": {}, "VA": {}, "VE": {}, "VG": {}, "VN": {}, "VU": {}, "YE": {},
	"ZA": {},
}

func getCountryCodes() []string {
	codes := make([]string, 0, len(validCountries))
	for code := range validCountries {
		codes = append(codes, code)
	}
	sort.Strings(codes)
	return codes
}
