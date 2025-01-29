package xparser

import (
	"reflect"
	"testing"
)

// Test struct
type SatellitePosition struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Altitude  float64 `json:"altitude"`
	Timestamp string  `json:"timestamp"`
	UID       string  `json:"uid"`
}

// SatelliteMapper function for testing (example: prefix name with "Mapped-")
func SatelliteMapper(input SatellitePosition) (SatellitePosition, error) {
	return SatellitePosition{
		ID:        input.ID,
		Name:      "Mapped-" + input.Name,
		Latitude:  input.Latitude,
		Longitude: input.Longitude,
		Altitude:  input.Altitude,
		Timestamp: input.Timestamp,
		UID:       input.UID,
	}, nil
}

// TestSerializeStruct tests serialization with and without a mapping function
func TestSerializeStruct(t *testing.T) {
	sat := SatellitePosition{
		ID:        "123",
		Name:      "Hubble",
		Latitude:  37.7749,
		Longitude: -122.4194,
		Altitude:  500.5,
		Timestamp: "2024-01-29T12:00:00Z",
		UID:       "user123",
	}

	// Serialize without mapping function
	data, err := SerializeStruct(sat, nil)
	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}

	// Check if all fields exist in the map
	expectedKeys := []string{"id", "name", "latitude", "longitude", "altitude", "timestamp", "uid"}
	for _, key := range expectedKeys {
		if _, exists := data[key]; !exists {
			t.Errorf("Expected key %s in serialized map, but it was missing", key)
		}
	}

	// Serialize with mapping function
	dataMapped, err := SerializeStruct(sat, SatelliteMapper)
	if err != nil {
		t.Fatalf("Serialization with mapper failed: %v", err)
	}

	// Check if the name was correctly modified
	if dataMapped["name"] != "Mapped-Hubble" {
		t.Errorf("Expected mapped name 'Mapped-Hubble', got %s", dataMapped["name"])
	}
}

// TestDeserializeStruct tests deserialization with and without a mapping function
func TestDeserializeStruct(t *testing.T) {
	data := map[string]interface{}{
		"id":        "123",
		"name":      "Hubble",
		"latitude":  37.7749,
		"longitude": -122.4194,
		"altitude":  500.5,
		"timestamp": "2024-01-29T12:00:00Z",
		"uid":       "user123",
	}

	// Deserialize without mapping function
	sat, err := DeserializeStruct[SatellitePosition](data, nil)
	if err != nil {
		t.Fatalf("Deserialization failed: %v", err)
	}

	// Expected struct
	expected := SatellitePosition{
		ID:        "123",
		Name:      "Hubble",
		Latitude:  37.7749,
		Longitude: -122.4194,
		Altitude:  500.5,
		Timestamp: "2024-01-29T12:00:00Z",
		UID:       "user123",
	}

	// Compare structs
	if !reflect.DeepEqual(sat, expected) {
		t.Errorf("Deserialized struct does not match expected struct.\nGot: %+v\nExpected: %+v", sat, expected)
	}

	// Deserialize with mapping function
	satMapped, err := DeserializeStruct[SatellitePosition](data, SatelliteMapper)
	if err != nil {
		t.Fatalf("Deserialization with mapper failed: %v", err)
	}

	// Check if the name was correctly modified
	if satMapped.Name != "Mapped-Hubble" {
		t.Errorf("Expected mapped name 'Mapped-Hubble', got %s", satMapped.Name)
	}
}
