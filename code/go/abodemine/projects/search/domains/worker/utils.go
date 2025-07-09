package worker

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// getFullStateName converts a state abbreviation to its full name
func getFullStateName(stateAbbr *string) string {
	if stateAbbr == nil {
		return ""
	}

	stateMap := map[string]string{
		"AL": "Alabama",
		"AK": "Alaska",
		"AZ": "Arizona",
		"AR": "Arkansas",
		"CA": "California",
		"CO": "Colorado",
		"CT": "Connecticut",
		"DE": "Delaware",
		"FL": "Florida",
		"GA": "Georgia",
		"HI": "Hawaii",
		"ID": "Idaho",
		"IL": "Illinois",
		"IN": "Indiana",
		"IA": "Iowa",
		"KS": "Kansas",
		"KY": "Kentucky",
		"LA": "Louisiana",
		"ME": "Maine",
		"MD": "Maryland",
		"MA": "Massachusetts",
		"MI": "Michigan",
		"MN": "Minnesota",
		"MS": "Mississippi",
		"MO": "Missouri",
		"MT": "Montana",
		"NE": "Nebraska",
		"NV": "Nevada",
		"NH": "New Hampshire",
		"NJ": "New Jersey",
		"NM": "New Mexico",
		"NY": "New York",
		"NC": "North Carolina",
		"ND": "North Dakota",
		"OH": "Ohio",
		"OK": "Oklahoma",
		"OR": "Oregon",
		"PA": "Pennsylvania",
		"RI": "Rhode Island",
		"SC": "South Carolina",
		"SD": "South Dakota",
		"TN": "Tennessee",
		"TX": "Texas",
		"UT": "Utah",
		"VT": "Vermont",
		"VA": "Virginia",
		"WA": "Washington",
		"WV": "West Virginia",
		"WI": "Wisconsin",
		"WY": "Wyoming",
		// Territories
		"AS": "American Samoa",
		"DC": "District of Columbia",
		"FM": "Federated States of Micronesia",
		"GU": "Guam",
		"MH": "Marshall Islands",
		"MP": "Northern Mariana Islands",
		"PW": "Palau",
		"PR": "Puerto Rico",
		"VI": "Virgin Islands",
		// Armed Forces (AE includes Europe, Africa, Canada, and the Middle East)
		"AA": "Armed Forces Americas",
		"AE": "Armed Forces Europe",
		"AP": "Armed Forces Pacific",
	}

	// make caps
	stateAbbrCaps := strings.ToUpper(*stateAbbr)

	if fullName, ok := stateMap[stateAbbrCaps]; ok {
		return fullName
	}

	return stateAbbrCaps // Return the abbreviation if no mapping found
}

// FormatAddressText formats address text to proper case (Title case with special handling for address-specific terms)
// It handles common address prefixes, suffixes, and special cases like "McKee" or "O'Brien"
func FormatAddressText(input string) string {
	if input == "" {
		return ""
	}

	// Handle all uppercase input
	if strings.ToUpper(input) == input {
		input = strings.ToLower(input)
	}

	// Create a title caser
	titleCaser := cases.Title(language.English)

	// Split the input into words
	words := strings.Fields(input)
	for i, word := range words {
		// Skip empty words
		if word == "" {
			continue
		}

		// Handle special cases for address components
		lowerWord := strings.ToLower(word)

		switch {
		// Directionals - keep lowercase when they're not the first word
		case lowerWord == "n" || lowerWord == "s" || lowerWord == "e" || lowerWord == "w" ||
			lowerWord == "ne" || lowerWord == "nw" || lowerWord == "se" || lowerWord == "sw" ||
			lowerWord == "north" || lowerWord == "south" || lowerWord == "east" || lowerWord == "west" ||
			lowerWord == "northeast" || lowerWord == "northwest" || lowerWord == "southeast" || lowerWord == "southwest":
			if i > 0 {
				words[i] = lowerWord
			} else {
				words[i] = titleCaser.String(lowerWord)
			}

		// Common address prefixes to keep lowercase
		case lowerWord == "de" || lowerWord == "del" || lowerWord == "la" || lowerWord == "las" ||
			lowerWord == "los" || lowerWord == "von" || lowerWord == "van" || lowerWord == "der":
			words[i] = lowerWord

		// Special case for Mc names (like McKee)
		case strings.HasPrefix(lowerWord, "mc"):
			if len(word) > 2 {
				words[i] = "Mc" + titleCaser.String(lowerWord[2:])
			} else {
				words[i] = "Mc"
			}

		// Special case for Mac names (like MacArthur)
		case strings.HasPrefix(lowerWord, "mac"):
			if len(word) > 3 {
				words[i] = "Mac" + titleCaser.String(lowerWord[3:])
			} else {
				words[i] = "Mac"
			}

		// Special case for O' names (like O'Brien)
		case strings.Contains(lowerWord, "o'"):
			parts := strings.Split(lowerWord, "o'")
			if len(parts) > 1 {
				words[i] = "O'" + titleCaser.String(parts[1])
			} else {
				words[i] = titleCaser.String(lowerWord)
			}

		// Default case - title case
		default:
			words[i] = titleCaser.String(lowerWord)
		}
	}

	return strings.Join(words, " ")
}

// CustomWriter wraps zerolog.ConsoleWriter to format numbers and durations
type CustomWriter struct {
	w zerolog.ConsoleWriter
}

func (w *CustomWriter) Write(p []byte) (n int, err error) {
	var event map[string]interface{}
	if err := json.Unmarshal(p, &event); err != nil {
		return w.w.Write(p) // If we can't parse as JSON, write original
	}

	// Format numeric values in the event
	for k, v := range event {
		switch val := v.(type) {
		case float64:
			if k == "duration" || strings.HasSuffix(k, "_duration") {
				event[k] = fmt.Sprintf("%dms", int(val))
			} else if val == float64(int(val)) {
				event[k] = formatNumber(int(val))
			} else {
				event[k] = fmt.Sprintf("%.2f", val)
			}
		case string:
			if strings.Contains(val, "ms") {
				event[k] = val // Keep ms suffix
			}
		}
	}

	// Convert back to JSON
	formatted, err := json.Marshal(event)
	if err != nil {
		return w.w.Write(p) // If marshal fails, write original
	}

	return w.w.Write(formatted)
}

// formatNumber formats a number with thousand separators
func formatNumber(n int) string {
	in := strconv.Itoa(n)
	numOfDigits := len(in)
	if n < 0 {
		numOfDigits-- // First character is the - sign (not a digit)
	}
	numOfCommas := (numOfDigits - 1) / 3

	out := make([]byte, len(in)+numOfCommas)
	if n < 0 {
		in = in[1:]
		out[0] = '-'
	}

	for i, j, k := len(in)-1, len(out)-1, 0; ; i, j = i-1, j-1 {
		out[j] = in[i]
		if i == 0 {
			break
		}
		k++
		if k == 3 {
			j--
			out[j] = ','
			k = 0
		}
	}
	return string(out)
}
