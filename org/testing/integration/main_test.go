package integration

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/cucumber/godog"
	"github.com/org/2112-space-lab/org/testing/integration/resources"
	"github.com/org/2112-space-lab/org/testing/pkg/fx"
	xtestartifact "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-artifact"
	xtestcommon "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-common"
	"github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-common/models"
	models_cont "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-container/models"
	models_db "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-db/models"
	xtestlog "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-log"
)

func TestMain(m *testing.M) {
	runRandID := xtestcommon.GetOrInitRunRandID()

	runLogger := slog.With(
		slog.String("runRandID", string(runRandID)),
	)
	ctx := context.Background()
	parseFlags(runLogger)
	artifactsBasePath := CLI.FlagWrapper.Options.ArtifactsPath
	xtestartifact.InitArtifactsFolderOrPanic(ctx, runLogger, artifactsBasePath)
	runLogger, teardownLogger := xtestlog.InitRunLoggerOrPanic(ctx, runLogger, CLI.FlagWrapper.Log.Level, CLI.FlagWrapper.Log.Format, runRandID, artifactsBasePath)
	defer func() {
		err := teardownLogger(ctx)
		if err != nil {
			slog.Error("failed to teardownLogger",
				slog.Any("error", err),
			)
		}
	}()

	teardown, err := resources.InitGlobalResourceManager(
		ctx,
		runLogger,
		GetDatabaseConnectionInfo(),
		CLI.FlagWrapper.Options.AppWorkspacePath,
		models_db.DatabaseMigrationPath(CLI.FlagWrapper.Options.AppMigrationsPath),
		models_cont.DockerContainerImage(CLI.FlagWrapper.Options.AppServiceDockerImage),
		models_cont.DockerContainerImage(CLI.FlagWrapper.Options.PropagatorDockerImage),
		models_cont.NetworkName(CLI.FlagWrapper.Options.DockerNetwork),
		models.TestRunningEnv(CLI.FlagWrapper.Options.TestRunningEnv),
	)
	if err != nil {
		slog.Error("failed to InitGlobalResourceManager",
			slog.Any("error", err),
		)
		runLogger.Error("failed to InitGlobalResourceManager",
			slog.Any("error", err),
		)
		os.Exit(1)
	}
	if teardown != nil {
		defer func() {
			err = teardown(ctx)
			if err != nil {
				slog.Error("failed to teardown resources",
					slog.Any("error", err),
				)
			}
		}()
	}
	status := m.Run()
	os.Exit(status)
}

func performAppTestSuite(t *testing.T, specSubFolder string) {
	o := opts
	o.Paths = []string{filepath.Join("specs/", specSubFolder)}
	o.TestingT = t
	suiteFolder := filepath.Join("./_artifacts/suites", specSubFolder)
	o.Format = fmt.Sprintf(
		"cucumber:%s/godog-cucumber.log.json,pretty:%s/godog-pretty.log,pretty",
		suiteFolder,
		suiteFolder,
	)
	suiteInit, scenarioInit := PrepareGodogInitializers(context.Background(), xtestlog.GetRunLogger(), specSubFolder, t.Name())
	suite := godog.TestSuite{
		Name:                 t.Name(),
		Options:              &o,
		TestSuiteInitializer: suiteInit,
		ScenarioInitializer:  scenarioInit,
	}
	status := suite.Run()

	if status == 2 {
		t.SkipNow()
	}

	if status != 0 {
		t.Fatalf("zero status code expected, %d received", status)
	}
}

func parseFlags(logger *slog.Logger) {
	kongAppTestArgs := fx.FilterSlice(os.Args[1:], keepKongFlags)
	godogArgs := fx.FilterSlice(os.Args, godogFlags)

	logger.Info("test run args",
		slog.Any("os.Args", os.Args),
		slog.Any("kongAppTestArgs", kongAppTestArgs),
		slog.Any("godogArgs", godogArgs),
	)

	parseAppTestKongFlags(logger, kongAppTestArgs)
	parseGodogFlags(logger, CLI.FlagWrapper.Options.ArtifactsPath, godogArgs)
	if CLI.FlagWrapper.ShowHelp {
		os.Exit(0)
	}
}
