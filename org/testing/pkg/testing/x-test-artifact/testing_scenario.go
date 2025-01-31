package xtestartifact

import (
	"context"
	"log/slog"
	"os"

	models_common "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-common/models"
)

func InitArtifactScenarioFolderOrPanic(ctx context.Context, logger *slog.Logger, suiteScenariosBasePath string, scenarioRandID models_common.ScenarioRandID) string {
	basePathStat, err := os.Stat(suiteScenariosBasePath)
	if err != nil {
		panic(err)
	}
	scenarioFolderPath := createArtifactSubfolderOrPanic(logger, suiteScenariosBasePath, string(scenarioRandID), basePathStat)
	return scenarioFolderPath
}
