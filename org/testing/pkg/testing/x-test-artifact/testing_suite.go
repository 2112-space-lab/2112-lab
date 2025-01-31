package xtestartifact

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
)

const (
	suiteSubfolderNameScenarios string = "scenarios"
)

func InitArtifactSuiteFolderOrPanic(ctx context.Context, logger *slog.Logger, artifactsBasePath string, suiteFolderName string) (string, string) {
	basePathStat, err := os.Stat(artifactsBasePath)
	if err != nil {
		panic(err)
	}
	suitesBasePath := filepath.Join(artifactsBasePath, artifactsSubfolderNameSuites)
	suitePath := createArtifactSubfolderOrPanic(logger, suitesBasePath, suiteFolderName, basePathStat)
	suiteScenariosBasePath := createArtifactSubfolderOrPanic(logger, suitePath, suiteSubfolderNameScenarios, basePathStat)
	return suitePath, suiteScenariosBasePath
}
