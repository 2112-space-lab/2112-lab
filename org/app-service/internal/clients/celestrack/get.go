package celestrack

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/org/2112-space-lab/org/app-service/internal/api/mappers"
	"github.com/org/2112-space-lab/org/app-service/internal/config"
	"github.com/org/2112-space-lab/org/go-utils/pkg/fx/xspace"
)

// CelestrackClient definition
type CelestrackClient struct {
	env *config.SEnv
}

// NewCelestrackClient constructor
func NewCelestrackClient(env *config.SEnv) *CelestrackClient {
	return &CelestrackClient{
		env: env,
	}
}

// FetchTLEFromSatCatByCategory fetches TLE from satcat catalog
func (client *CelestrackClient) FetchTLEFromSatCatByCategory(ctx context.Context, category string) ([]*mappers.RawTLE, error) {
	if category == "" {
		return nil, fmt.Errorf("category is required")
	}

	baseUrl := client.env.EnvVars.Celestrack.BaseUrl
	// Construct the URL for the category
	url := fmt.Sprintf("%s?GROUP=%s", baseUrl, category)

	// Fetch the TLE data
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch TLE data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch TLE data: HTTP status %d", resp.StatusCode)
	}

	// Parse the TLE data
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read TLE data: %v", err)
	}

	lines := bytes.Split(body, []byte("\n"))
	var tles []*mappers.RawTLE
	for i := 0; i < len(lines)-1; i += 3 {
		if len(lines[i]) == 0 || len(lines[i+1]) == 0 || len(lines[i+2]) == 0 {
			continue
		}

		// Extract NORAD ID from Line 1 (positions 3–7 as per TLE format)
		line1 := strings.TrimSpace(string(lines[i+1]))
		if len(line1) < 7 {
			continue // Skip invalid lines
		}
		noradID := strings.TrimSpace(line1[2:7]) // Extract NORAD ID

		line2 := strings.TrimSpace(string(lines[i+2]))

		tles = append(tles, &mappers.RawTLE{
			NoradID: noradID,
			Line1:   line1,
			Line2:   line2,
		})
	}

	return tles, nil
}

// FetchSatelliteMetadata fetches metadata for satellites from CelesTrak's SATCAT.
func (client *CelestrackClient) FetchSatelliteMetadata(ctx context.Context) ([]*mappers.SatelliteMetadata, error) {
	// Create an HTTP request with the provided context
	satcatUrl := client.env.EnvVars.Celestrack.Satcat
	req, err := http.NewRequestWithContext(ctx, "GET", satcatUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Execute the HTTP request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch SATCAT data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch SATCAT data: HTTP status %d", resp.StatusCode)
	}

	// Parse the CSV data
	reader := csv.NewReader(resp.Body)
	reader.FieldsPerRecord = -1 // Allow variable number of fields per record

	// Read all records
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read SATCAT data: %v", err)
	}

	// Extract header indices
	header := records[0]
	indices := make(map[string]int)
	for i, field := range header {
		indices[strings.TrimSpace(field)] = i
	}

	// Parse records into SatelliteMetadata
	var satellites []*mappers.SatelliteMetadata
	for _, record := range records[1:] {
		launchDate, err := time.Parse("2006-01-02", record[indices["LAUNCH_DATE"]])
		if err != nil {
			continue // Skip records with invalid launch dates
		}

		var decayDate *time.Time
		if dateStr := record[indices["DECAY_DATE"]]; dateStr != "" {
			d, err := time.Parse("2006-01-02", dateStr)
			if err == nil {
				decayDate = &d
			}
		}

		var apogee, perigee float64
		if apogeeStr := record[indices["APOGEE"]]; apogeeStr != "" {
			apogee, _ = parseFloat(apogeeStr)
		}
		if perigeeStr := record[indices["PERIGEE"]]; perigeeStr != "" {
			perigee, _ = parseFloat(perigeeStr)
		}

		altitude := xspace.ComputeAverageAltitude(apogee, perigee)

		satellite := &mappers.SatelliteMetadata{
			NoradID:        record[indices["NORAD_CAT_ID"]],
			Name:           record[indices["OBJECT_NAME"]],
			IntlDesignator: record[indices["OBJECT_ID"]],
			LaunchDate:     launchDate,
			DecayDate:      decayDate,
			ObjectType:     record[indices["OBJECT_TYPE"]],
			Owner:          record[indices["OWNER"]],
			Altitude:       &altitude,
		}

		// Optional fields
		if periodStr := record[indices["PERIOD"]]; periodStr != "" {
			if period, err := parseFloat(periodStr); err == nil {
				satellite.Period = &period
			}
		}
		if inclinationStr := record[indices["INCLINATION"]]; inclinationStr != "" {
			if inclination, err := parseFloat(inclinationStr); err == nil {
				satellite.Inclination = &inclination
			}
		}
		if rcsStr := record[indices["RCS"]]; rcsStr != "" {
			if rcs, err := parseFloat(rcsStr); err == nil {
				satellite.RCS = &rcs
			}
		}

		satellites = append(satellites, satellite)
	}

	return satellites, nil
}

// parseFloat is a helper function to parse a float64 from a string.
func parseFloat(s string) (float64, error) {
	return strconv.ParseFloat(strings.TrimSpace(s), 64)
}
