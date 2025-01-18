package utils

import (
	"fmt"
	"strings"
)

// GetStateCode returns the code for a given Indian state.
func GetStateCode(state string) (string, error) {
	stateCodes := map[string]string{
		"Andhra Pradesh":              "AP",
		"Arunachal Pradesh":           "AR",
		"Assam":                       "AS",
		"Bihar":                       "BR",
		"Chhattisgarh":                "CG",
		"Goa":                         "GA",
		"Gujarat":                     "GJ",
		"Haryana":                     "HR",
		"Himachal Pradesh":            "HP",
		"Jharkhand":                   "JH",
		"Karnataka":                   "KA",
		"Kerala":                      "KL",
		"Madhya Pradesh":              "MP",
		"Maharashtra":                 "MH",
		"Manipur":                     "MN",
		"Meghalaya":                   "ML",
		"Mizoram":                     "MZ",
		"Nagaland":                    "NL",
		"Odisha":                      "OD",
		"Punjab":                      "PB",
		"Rajasthan":                   "RJ",
		"Sikkim":                      "SK",
		"Tamil Nadu":                  "TN",
		"Telangana":                   "TS",
		"Tripura":                     "TR",
		"Uttar Pradesh":               "UP",
		"Uttarakhand":                 "UK",
		"West Bengal":                 "WB",
		"Andaman and Nicobar Islands": "AN",
		"Chandigarh":                  "CH",
		"Dadra and Nagar Haveli and Daman and Diu": "DN",
		"Delhi":             "DL",
		"Lakshadweep":       "LD",
		"Puducherry":        "PY",
		"Ladakh":            "LA",
		"Jammu and Kashmir": "JK",
	}

	if code, exists := stateCodes[state]; exists {
		return code, nil
	}

	return "", fmt.Errorf("state code not found for %s", state)
}

// GetCityCode returns a three-letter code for a given Indian city.
func GetCityCode(city string) (string, error) {
	cityCodes := map[string]string{
		"New Delhi":          "NDL",
		"Mumbai":             "MUM",
		"Bangalore":          "BLR",
		"Hyderabad":          "HYD",
		"Chennai":            "CHE",
		"Kolkata":            "KOL",
		"Pune":               "PUN",
		"Ahmedabad":          "AMD",
		"Jaipur":             "JAI",
		"Lucknow":            "LKO",
		"Chandigarh":         "CHD",
		"Patna":              "PAT",
		"Bhopal":             "BPL",
		"Thiruvananthapuram": "TVM",
		"Vadodara":           "VAD",
		"Nagpur":             "NAG",
		"Indore":             "IND",
		"Visakhapatnam":      "VIZ",
		"Surat":              "SUR",
		"Coimbatore":         "CBE",
		"Kanpur":             "KNP",
		"Agra":               "AGR",
		"Varanasi":           "VNS",
		"Ranchi":             "RNC",
		"Jamshedpur":         "JSP",
		"Dehradun":           "DED",
		"Shimla":             "SHL",
		"Srinagar":           "SXR",
		"Guwahati":           "GUW",
		"Shillong":           "SHG",
		"Imphal":             "IMP",
		"Aizawl":             "AIZ",
		"Kohima":             "KOH",
		"Itanagar":           "ITA",
		"Panaji":             "PNJ",
		"Raipur":             "RPR",
		"Amritsar":           "ASR",
		"Ludhiana":           "LDH",
		"Faridabad":          "FRD",
		"Ghaziabad":          "GZB",
		"Noida":              "NOD",
		"Gurgaon":            "GGN",
		"Madurai":            "MDU",
		"Trichy":             "TRZ",
		"Salem":              "SLM",
		"Jodhpur":            "JOD",
		"Udaipur":            "UDA",
		"Kota":               "KOT",
		"Meerut":             "MRT",
		"Aligarh":            "ALG",
		"Gwalior":            "GWL",
		"Jabalpur":           "JBP",
		"Ajmer":              "AJM",
		"Nashik":             "NSK",
		"Kolhapur":           "KOP",
		"Solapur":            "SLP",
		"Hubli":              "HBL",
		"Belgaum":            "BEL",
		"Mysore":             "MYS",
		"Rajkot":             "RAJ",
		"Jamnagar":           "JAM",
		"Bhavnagar":          "BHV",
		"Anand":              "AND",
		"Bhuj":               "BHJ",
		"Porbandar":          "PBD",
		"Jalandhar":          "JAL",
		"Ambala":             "AMB",
		"Bilaspur":           "BSP",
		"Durgapur":           "DGP",
		"Asansol":            "ASN",
		"Howrah":             "HWH",
	}

	// Normalize input city name to handle case sensitivity and trim spaces
	if code, exists := cityCodes[strings.TrimSpace(city)]; exists {
		return code, nil
	}

	return "", fmt.Errorf("city code not found for %s", city)
}
