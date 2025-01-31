package xtestartifact

import (
	"context"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
)

const (
	artifactsSubfolderNameSuites  string = "suites"
	artifactsSubfolderNameResults string = "results"
	artifactsSubfolderNameReports string = "reports"
)

func InitArtifactsFolderOrPanic(ctx context.Context, logger *slog.Logger, artifactsBasePath string) {
	logger = logger.With(
		slog.String("artifactsPath", artifactsBasePath),
	)
	logger.Info("initializing test artifacts output folder")
	basePathStat, err := os.Stat(".")
	if err != nil {
		logger.Error("failed to stat current directory",
			slog.Any("error", err),
		)
		panic(err)
	}

	err = os.RemoveAll(artifactsBasePath)
	if err != nil {
		logger.Error("failed to clear existing old test artifacts output folder",
			slog.Any("error", err),
		)
		panic(err)
	}
	logger.Info("cleared test artifacts output folder")

	err = os.MkdirAll(artifactsBasePath, basePathStat.Mode())
	if err != nil {
		logger.Error("failed to (re)create test artifacts output folder",
			slog.Any("error", err),
		)
		panic(err)
	}
	logger.Info("(re)created test artifacts output folder")

	_ = createArtifactSubfolderOrPanic(logger, artifactsBasePath, artifactsSubfolderNameResults, basePathStat)
	_ = createArtifactSubfolderOrPanic(logger, artifactsBasePath, artifactsSubfolderNameReports, basePathStat)
	_ = createArtifactSubfolderOrPanic(logger, artifactsBasePath, artifactsSubfolderNameSuites, basePathStat)
}

func createArtifactSubfolderOrPanic(logger *slog.Logger, artifactsPath string, subFolderName string, basePathStat fs.FileInfo) string {
	logger = logger.With(
		slog.String("subFolderName", subFolderName),
		slog.String("basePath", artifactsPath),
	)
	newPath := filepath.Join(artifactsPath, subFolderName)
	err := os.MkdirAll(newPath, basePathStat.Mode())
	if err != nil {
		logger.Error("failed to create test artifacts sub-folder",
			slog.Any("error", err),
		)
		panic(err)
	}
	logger.Log(context.Background(), slog.LevelDebug.Level(), "created test artifacts sub-folder")
	// logger.Log(context.Background(), xlog.LevelTrace, "created test artifacts sub-folder")
	return newPath
}
