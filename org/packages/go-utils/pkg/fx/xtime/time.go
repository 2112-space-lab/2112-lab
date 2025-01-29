package xtime

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// TLEFormat converts a TLE epoch string into a `time.Time` object.
// The TLE epoch is typically in the format YYDDD.DDDDDDDD (e.g., 22365.12345678).
func TLEFormat(epochStr string) (time.Time, error) {
	// Ensure epochStr is the expected minimum length (YYDDD)
	if len(epochStr) < 5 {
		return time.Time{}, fmt.Errorf("invalid TLE epoch string: %s", epochStr)
	}

	// Split into year and day parts
	yearStr := epochStr[:2]    // First two characters: last two digits of the year
	dayPartStr := epochStr[2:] // Remaining characters: day of year and fractional day

	// Convert year part to an integer and adjust for the full year
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid year in TLE epoch: %v", err)
	}
	// Handle Y2K problem: years 00-56 are 2000-2056, 57-99 are 1957-1999
	if year < 57 {
		year += 2000
	} else {
		year += 1900
	}

	// Split dayPartStr into whole days and fractional parts
	var dayOfYear int
	var fractionOfDay float64
	if len(dayPartStr) > 3 {
		dayOfYear, err = strconv.Atoi(dayPartStr[:3]) // Extract whole days (DDD)
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid day in TLE epoch: %v", err)
		}
		fractionOfDay, err = strconv.ParseFloat("0"+dayPartStr[3:], 64) // Extract fractional part
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid fractional day in TLE epoch: %v", err)
		}
	} else {
		return time.Time{}, fmt.Errorf("day and fractional part too short in TLE epoch: %s", dayPartStr)
	}

	// Calculate the time of the epoch
	startOfYear := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	parsedTime := startOfYear.AddDate(0, 0, dayOfYear-1).Add(time.Duration(fractionOfDay * 24 * float64(time.Hour)))

	return parsedTime, nil
}

// FromRawTLE attempts to parse a time from a given TLE epoch string.
func FromRawTLE(epochStr string) (time.Time, error) {
	return TLEFormat(strings.TrimSpace(epochStr))
}

func ParseEpoch(line1 string) (time.Time, error) {
	// Extract epoch substring from the TLE line
	if len(line1) < 32 {
		return time.Time{}, fmt.Errorf("invalid TLE line: epoch data missing")
	}
	epochStr := strings.TrimSpace(line1[18:32])
	return FromRawTLE(epochStr)
}
