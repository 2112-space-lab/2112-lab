package models

type RunRandID string
type SuiteName string
type ScenarioRandID string

type ScenarioInfo struct {
	RunRandID      RunRandID
	SuiteName      SuiteName
	ScenarioRandID ScenarioRandID
}

func NewScenarioInfo(
	runRandID RunRandID,
	suiteName SuiteName,
	scenarioRandID ScenarioRandID,
) ScenarioInfo {
	return ScenarioInfo{
		RunRandID:      runRandID,
		SuiteName:      suiteName,
		ScenarioRandID: scenarioRandID,
	}
}
