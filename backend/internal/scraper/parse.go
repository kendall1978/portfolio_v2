package scraper

import (
	"strconv"
	"strings"
)

// StateNames maps 2-letter postal codes to full state names.
var StateNames = map[string]string{
	"AL": "Alabama", "AK": "Alaska", "AZ": "Arizona", "AR": "Arkansas",
	"CA": "California", "CO": "Colorado", "CT": "Connecticut", "DE": "Delaware",
	"DC": "District of Columbia", "FL": "Florida", "GA": "Georgia", "HI": "Hawaii",
	"ID": "Idaho", "IL": "Illinois", "IN": "Indiana", "IA": "Iowa",
	"KS": "Kansas", "KY": "Kentucky", "LA": "Louisiana", "ME": "Maine",
	"MD": "Maryland", "MA": "Massachusetts", "MI": "Michigan", "MN": "Minnesota",
	"MS": "Mississippi", "MO": "Missouri", "MT": "Montana", "NE": "Nebraska",
	"NV": "Nevada", "NH": "New Hampshire", "NJ": "New Jersey", "NM": "New Mexico",
	"NY": "New York", "NC": "North Carolina", "ND": "North Dakota", "OH": "Ohio",
	"OK": "Oklahoma", "OR": "Oregon", "PA": "Pennsylvania", "RI": "Rhode Island",
	"SC": "South Carolina", "SD": "South Dakota", "TN": "Tennessee", "TX": "Texas",
	"UT": "Utah", "VT": "Vermont", "VA": "Virginia", "WA": "Washington",
	"WV": "West Virginia", "WI": "Wisconsin", "WY": "Wyoming",
	"GU": "Guam", "PR": "Puerto Rico", "VI": "Virgin Islands",
}

// StateName returns the full state name for a 2-letter postal code.
func StateName(code string) string {
	return StateNames[code]
}

// ParseNumber strips commas and converts a string to float64.
// Returns (0, false) for suppressed values (*, #, empty).
func ParseNumber(s string) (float64, bool) {
	s = strings.TrimSpace(s)
	if s == "" || s == "*" || s == "#" || s == "**" {
		return 0, false
	}
	s = strings.ReplaceAll(s, ",", "")
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, false
	}
	return v, true
}

// ParseInt strips commas and converts a string to int.
// Returns (0, false) for suppressed values.
func ParseInt(s string) (int, bool) {
	v, ok := ParseNumber(s)
	if !ok {
		return 0, false
	}
	return int(v), true
}

// ShouldSkipOccCode returns true for group totals (codes ending in -0000).
func ShouldSkipOccCode(code string) bool {
	return len(code) >= 7 && code[len(code)-4:] == "0000"
}