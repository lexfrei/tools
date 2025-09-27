package main

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/cockroachdb/errors"
)

// Regular expression patterns for parsing different value formats.
var (
	// durationHoursPattern matches time in HH:MM:SS format.
	durationHoursPattern = regexp.MustCompile(`^(-?\d+,?\d*?):(\d+):(\d+)$`)
	// durationMinutesPattern matches time in MM:SS format.
	durationMinutesPattern = regexp.MustCompile(`^(-?\d+):(\d+)$`)
	// intPattern matches integer values with optional commas and percentage signs.
	intPattern = regexp.MustCompile(`^-?\d+(,\d+)*%?$`)
	// floatPattern matches float values with optional commas.
	floatPattern = regexp.MustCompile(`^-?\d+(,\d+)*\.\d+$`)
)

// ParseValue converts various string representations into appropriate Go types.
// It handles duration formats (HH:MM:SS, MM:SS) converting them to seconds,
// percentages, integers with commas, and floats.
// Returns 0 for special values like "--" or "NaN".
func ParseValue(input string) interface{} {
	if input == "" {
		return ""
	}

	// Handle special cases for missing or invalid data.
	if input == "--" || input == "NaN" {
		return 0
	}

	// Try to parse as duration first.
	value, err := parseDuration(input)
	if err == nil {
		return value
	}

	// Try to parse as numeric values.
	numValue, numErr := parseNumeric(input)
	if numErr == nil {
		return numValue
	}

	// Return original string if no pattern matches.
	return input
}

// parseDuration attempts to parse duration strings in HH:MM:SS or MM:SS format.
func parseDuration(input string) (int, error) {
	// Duration format in hour:min:sec => seconds.
	if matches := durationHoursPattern.FindStringSubmatch(input); matches != nil {
		return parseDurationHours(matches)
	}

	// Duration format in min:sec => seconds.
	if matches := durationMinutesPattern.FindStringSubmatch(input); matches != nil {
		return parseDurationMinutes(matches)
	}

	return 0, errors.New("not a duration format")
}

// parseDurationHours parses HH:MM:SS format and returns total seconds.
func parseDurationHours(matches []string) (int, error) {
	hours, err := strconv.Atoi(strings.ReplaceAll(matches[1], ",", ""))
	if err != nil {
		return 0, errors.Wrap(err, "failed to parse hours")
	}

	minutes, err := strconv.Atoi(matches[2])
	if err != nil {
		return 0, errors.Wrap(err, "failed to parse minutes")
	}

	seconds, err := strconv.Atoi(matches[3])
	if err != nil {
		return 0, errors.Wrap(err, "failed to parse seconds")
	}

	const (
		secondsPerHour   = 3600
		secondsPerMinute = 60
	)

	return hours*secondsPerHour + minutes*secondsPerMinute + seconds, nil
}

// parseDurationMinutes parses MM:SS format and returns total seconds.
func parseDurationMinutes(matches []string) (int, error) {
	minutes, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, errors.Wrap(err, "failed to parse minutes")
	}

	seconds, err := strconv.Atoi(matches[2])
	if err != nil {
		return 0, errors.Wrap(err, "failed to parse seconds")
	}

	const secondsPerMinute = 60

	return minutes*secondsPerMinute + seconds, nil
}

// parseNumeric attempts to parse numeric values (integers and floats).
func parseNumeric(input string) (interface{}, error) {
	// Integer format (including percentages).
	if intPattern.MatchString(input) {
		cleanValue := strings.ReplaceAll(strings.ReplaceAll(input, "%", ""), ",", "")
		value, err := strconv.Atoi(cleanValue)
		if err != nil {
			return 0, errors.Wrap(err, "failed to parse integer")
		}

		return value, nil
	}

	// Float format.
	if floatPattern.MatchString(input) {
		cleanValue := strings.ReplaceAll(input, ",", "")
		value, err := strconv.ParseFloat(cleanValue, 64)
		if err != nil {
			return 0, errors.Wrap(err, "failed to parse float")
		}

		return value, nil
	}

	return nil, errors.New("not a numeric format")
}

// StringToSnakeCase converts a string to snake_case format.
// It removes accents, handles camelCase conversion, and replaces
// non-alphanumeric characters with underscores.
func StringToSnakeCase(input string) string {
	if input == "" {
		return ""
	}

	// Remove accents and normalize the string.
	cleaned := removeAccents(input)
	// Remove "- " sequences.
	cleaned = strings.ReplaceAll(cleaned, "- ", "")

	var result strings.Builder
	result.Grow(len(cleaned) * 2) // Pre-allocate to avoid reallocations.

	for i, char := range cleaned {
		if unicode.IsUpper(char) && i > 0 {
			// Check if previous character is lowercase.
			prevChar := rune(cleaned[i-1])
			if unicode.IsLower(prevChar) {
				result.WriteRune('_')
			}
		}

		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			result.WriteRune(unicode.ToLower(char))
		} else {
			result.WriteRune('_')
		}
	}

	// Clean up multiple underscores and trim.
	resultStr := result.String()
	resultStr = regexp.MustCompile(`_+`).ReplaceAllString(resultStr, "_")
	resultStr = strings.Trim(resultStr, "_")

	return resultStr
}

// removeAccents removes accents from Unicode characters.
// This is a simplified version that handles common accented characters.
func removeAccents(input string) string {
	// Map of common accented characters to their base forms.
	accentMap := map[rune]rune{
		'à': 'a', 'á': 'a', 'â': 'a', 'ã': 'a', 'ä': 'a', 'å': 'a',
		'è': 'e', 'é': 'e', 'ê': 'e', 'ë': 'e',
		'ì': 'i', 'í': 'i', 'î': 'i', 'ï': 'i',
		'ò': 'o', 'ó': 'o', 'ô': 'o', 'õ': 'o', 'ö': 'o',
		'ù': 'u', 'ú': 'u', 'û': 'u', 'ü': 'u',
		'ý': 'y', 'ÿ': 'y',
		'ñ': 'n',
		'ç': 'c',
		'À': 'A', 'Á': 'A', 'Â': 'A', 'Ã': 'A', 'Ä': 'A', 'Å': 'A',
		'È': 'E', 'É': 'E', 'Ê': 'E', 'Ë': 'E',
		'Ì': 'I', 'Í': 'I', 'Î': 'I', 'Ï': 'I',
		'Ò': 'O', 'Ó': 'O', 'Ô': 'O', 'Õ': 'O', 'Ö': 'O',
		'Ù': 'U', 'Ú': 'U', 'Û': 'U', 'Ü': 'U',
		'Ý': 'Y', 'Ÿ': 'Y',
		'Ñ': 'N',
		'Ç': 'C',
	}

	var result strings.Builder
	result.Grow(len(input))

	for _, char := range input {
		if replacement, exists := accentMap[char]; exists {
			result.WriteRune(replacement)
		} else {
			result.WriteRune(char)
		}
	}

	return result.String()
}

// ExtractStatsHeroClass extracts the "option-N" pattern from CSS class strings.
// This is used to identify hero-specific CSS classes in the DOM.
func ExtractStatsHeroClass(heroClass string) string {
	const optionPrefix = "option-"
	startIndex := strings.Index(heroClass, optionPrefix)
	if startIndex == -1 {
		return ""
	}

	endIndex := startIndex + len(optionPrefix)
	// Continue while we have digits.
	for endIndex < len(heroClass) && unicode.IsDigit(rune(heroClass[endIndex])) {
		endIndex++
	}

	return heroClass[startIndex:endIndex]
}
