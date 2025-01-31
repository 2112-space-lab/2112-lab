package xtestcommon

import "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-common/models"

func PrepareNewScenarioInfo(suiteName models.SuiteName) models.ScenarioInfo {
	info := models.NewScenarioInfo(
		GetOrInitRunRandID(),
		suiteName,
		GenerateScenarioRandID(),
	)
	return info
}
