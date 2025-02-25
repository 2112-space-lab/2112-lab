package state

import (
	"log/slog"

	stateservice "github.com/org/2112-space-lab/org/testing/pkg/testing/state-service"
	models_common "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-common/models"
	xteststate "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-state"
)

type ServiceScenarioState struct {
	*xteststate.ScenarioBaseState
	xteststate.TimeCheckpointState

	stateservice.AppDbState
	stateservice.AppState
	stateservice.PropagatorState
	stateservice.EventDetectorState
}

func RegisterCleanServiceScenarioState(info models_common.ScenarioInfo, logger *slog.Logger, scenarioFolderPath string) *ServiceScenarioState {
	base := xteststate.NewScenarioBaseState(info, logger, scenarioFolderPath)

	state := &ServiceScenarioState{
		ScenarioBaseState:   base,
		TimeCheckpointState: xteststate.NewTimeCheckpointState(),
		AppDbState:          stateservice.NewAppDBState(),
		AppState:            stateservice.NewAppState(),
		PropagatorState:     stateservice.NewPropagatorState(),
		EventDetectorState:  stateservice.NewEventDetectorState(),
	}
	return state
}
