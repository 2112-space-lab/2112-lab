package models

// GameContext represents a game context.
type GameContext struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    bool   `json:"isActive"`
}

// GameSatellitesContext represents a game context.
type GameSatellitesContext struct {
	Name           string   `json:"name"`
	SatelliteNames []string `json:"satelliteNames"`
}

type SatelliteDefinition struct {
	SatelliteName string `json:"SatelliteName"`
}
