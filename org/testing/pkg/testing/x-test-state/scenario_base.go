package xteststate

import (
	"log/slog"

	models_common "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-common/models"
)

type ScenarioBaseState struct {
	logger             *slog.Logger
	info               models_common.ScenarioInfo
	scenarioFolderPath string
}

func NewScenarioBaseState(info models_common.ScenarioInfo, logger *slog.Logger, scenarioFolderPath string) *ScenarioBaseState {
	state := &ScenarioBaseState{
		info:               info,
		scenarioFolderPath: scenarioFolderPath,
		logger:             logger,
	}
	return state
}

func (s *ScenarioBaseState) GetScenarioInfo() models_common.ScenarioInfo {
	return s.info
}

func (s *ScenarioBaseState) GetScenarioFolder() string {
	return s.scenarioFolderPath
}

func (s *ScenarioBaseState) GetLogger() *slog.Logger {
	if s.logger == nil {
		panic("unexpected nil scenario logger")
	}
	return s.logger
}
