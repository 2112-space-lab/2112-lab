package integration

import (
	"context"
	"log/slog"
	"time"

	"github.com/cucumber/godog"
	"github.com/org/2112-space-lab/org/testing/integration/resources"
	"github.com/org/2112-space-lab/org/testing/integration/state"
	"github.com/org/2112-space-lab/org/testing/integration/steps"
	"github.com/org/2112-space-lab/org/testing/pkg/fx"
	xtestartifact "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-artifact"
	xtestcommon "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-common"
	models_common "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-common/models"
	xtestcontainer "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-container"
	xtestlog "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-log"
)

func PrepareGodogInitializers(
	ctx context.Context,
	runLogger *slog.Logger,
	suiteFolderName string,
	suiteName string,
) (
	suiteInit func(tsc *godog.TestSuiteContext),
	scenarioInit func(sc *godog.ScenarioContext),
) {
	suiteScenarioBasePath, suiteLogger, suiteInit := prepareSuiteInitializer(ctx, runLogger, suiteFolderName, suiteName)
	scenarioInit = prepareScenarioInitializer(ctx, runLogger, suiteLogger, suiteName, suiteScenarioBasePath)
	return suiteInit, scenarioInit
}

func prepareSuiteInitializer(ctx context.Context, runLogger *slog.Logger, suiteFolderName string, suiteName string) (string, *slog.Logger, func(tsc *godog.TestSuiteContext)) {
	artifactsBasePath := CLI.FlagWrapper.Options.ArtifactsPath
	suiteBasePath, suiteScenariosPath := xtestartifact.InitArtifactSuiteFolderOrPanic(ctx, runLogger, artifactsBasePath, suiteFolderName)
	suiteLogger, teardownSuiteLogger := xtestlog.PrepareSuiteLogger(ctx, runLogger, CLI.FlagWrapper.Log.Level, CLI.FlagWrapper.Log.Format, xtestcommon.GetOrInitRunRandID(), suiteBasePath, suiteName)
	return suiteScenariosPath, suiteLogger, func(tsc *godog.TestSuiteContext) {
		startSuite := time.Now()
		tsc.BeforeSuite(func() {
			startSuite = time.Now()
			runLogger.Info("starting suite",
				slog.Group("suiteCtx",
					slog.String("suiteName", suiteName),
				),
			)
			suiteLogger.Info("starting suite")
		})
		tsc.AfterSuite(func() {
			defer func() {
				err := teardownSuiteLogger(ctx)
				if err != nil {
					slog.Error("failed to teardownSuiteLogger",
						slog.Any("error", err),
					)
				}
			}()
			elapsed := time.Since(startSuite)
			runLogger.Info("finished suite",
				slog.Group("suiteCtx",
					slog.String("suiteName", suiteName),
					slog.String("suiteDuration", elapsed.String()),
				),
			)
			suiteLogger.Info("finished suite",
				slog.String("suiteDuration", elapsed.String()),
			)

		})
	}
}

func prepareScenarioInitializer(baseCtx context.Context, runLogger *slog.Logger, suiteLogger *slog.Logger, suiteName string, suiteScenarioBasePath string) func(sc *godog.ScenarioContext) {
	return func(sc *godog.ScenarioContext) {
		scenarioInfo := xtestcommon.PrepareNewScenarioInfo(models_common.SuiteName(suiteName))
		scenarioFolderPath := xtestartifact.InitArtifactScenarioFolderOrPanic(baseCtx, suiteLogger, suiteScenarioBasePath, scenarioInfo.ScenarioRandID)
		scenarioLogger, teardownScenarioLogger := xtestlog.PrepareScenarioLogger(baseCtx, suiteLogger, CLI.FlagWrapper.Log.Level, CLI.FlagWrapper.Log.Format, scenarioFolderPath, scenarioInfo)

		resourceManager := resources.GetGlobalResourceManager()
		scenarioState := state.RegisterCleanServiceScenarioState(scenarioInfo, scenarioLogger, scenarioFolderPath)

		scenarioStart := time.Now()
		sc.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
			scenarioStart = time.Now()
			beforeScenarioLogging(runLogger, suiteLogger, scenarioLogger, sc, scenarioStart, scenarioInfo)
			return ctx, nil
		})
		sc.After(func(ctx context.Context, sc *godog.Scenario, scenarioErr error) (c context.Context, afterErr error) {
			defer func() {
				afterScenarioLogging(runLogger, suiteLogger, scenarioLogger, sc, scenarioErr, afterErr, scenarioStart, scenarioInfo)
				err := teardownScenarioLogger(ctx)
				if err != nil {
					slog.Error("failed to teardownScenarioLogger",
						slog.Any("error", err),
					)
				}
			}()

			errCompletion := scenarioState.WaitBackgroundsCompletionWithTimeout(15 * time.Second)
			backgroundErrors := scenarioState.GetBackgroundErrors()
			op := scenarioState.GetBackgroundOperations()
			var errBackground error
			if len(backgroundErrors) > 0 {
				scenarioLogger.Error("background steps summary with error",
					slog.Any("backgroundOperations", op),
					slog.Any("backgroundErrors", backgroundErrors),
				)
				errBackground = fx.FlattenErrorsIfAnyWithPath("backgroundSteps", backgroundErrors...)
			} else {
				scenarioLogger.Info("background steps summary",
					slog.Any("backgroundOperations", op),
				)

			}
			// teardown containers
			// scenarioState.CancelAllV2Stream()
			appContainers := fx.Values(scenarioState.GetScenarioAppServiceContainers())
			errAppTeardown := xtestcontainer.TeardownAllContainers(ctx, scenarioLogger, appContainers...)

			var errDB error
			if !CLI.FlagWrapper.Options.PreserveDatabases {
				appDBs := scenarioState.GetScenarioAppDatabases(ctx)
				errDB := resourceManager.RootDatabaseManager.DropAllDatabases(ctx, scenarioLogger, appDBs...)
				if errDB != nil {
					scenarioLogger.Error("failed to drop scenario databases",
						slog.Any(xtestlog.AttrErrorKey, errDB),
					)
				}
				scenarioLogger.Info("dropped scenario databases")
			}
			return ctx, fx.FlattenErrorsIfAny(errCompletion, errBackground, errAppTeardown, errDB)
		})

		stepStart := time.Now()
		sc.StepContext().Before(func(ctx context.Context, st *godog.Step) (context.Context, error) {
			stepStart = time.Now()
			scenarioLogger.Debug("before step", slog.Group("stepCtx",
				slog.String("stepText", st.Text),
			))
			return ctx, nil
		})
		sc.StepContext().After(func(ctx context.Context, st *godog.Step, status godog.StepResultStatus, err error) (context.Context, error) {
			elapsed := time.Since(stepStart)
			stepGroup := slog.Group("stepCtx",
				slog.String("stepText", st.Text),
				slog.String("stepDuration", elapsed.String()),
			)

			if err != nil {
				slog.Error("step failed",
					slog.Any(xtestlog.AttrErrorKey, err),
					stepGroup,
				)
				runLogger.Error("step failed",
					slog.Any(xtestlog.AttrErrorKey, err),
					stepGroup,
				)
				scenarioLogger.Error("step failed",
					slog.Any(xtestlog.AttrErrorKey, err),
					stepGroup,
				)
			} else {
				scenarioLogger.Info("step success", stepGroup)
			}
			return ctx, nil
		})

		steps.RegisterAppDatabaseSteps(sc, scenarioState, resourceManager.AppDatabaseManager)
		steps.RegisterAppServiceSteps(sc, scenarioState, resourceManager.ServiceContainerManager)
		steps.RegisterPropagatorServiceSteps(sc, scenarioState, resourceManager.ServiceContainerManager)
		steps.RegisterTimeCheckpointSteps(sc, scenarioState)
	}
}

func beforeScenarioLogging(runLogger, suiteLogger, scenarioLogger *slog.Logger, sc *godog.Scenario, scenarioStart time.Time, scenarioInfo models_common.ScenarioInfo) {
	msg := "before scenario"
	runLogger.Debug(msg, slog.Group("scenarioCtx",
		slog.String("suiteName", string(scenarioInfo.SuiteName)),
		slog.String("scenarioRandID", string(scenarioInfo.ScenarioRandID)),
		slog.String("scenarioUri", sc.Uri),
		slog.String("scenarioName", sc.Name),
	))
	suiteLogger.Debug(msg, slog.Group("scenarioCtx",
		slog.String("scenarioRandID", string(scenarioInfo.ScenarioRandID)),
		slog.String("scenarioUri", sc.Uri),
		slog.String("scenarioName", sc.Name),
	))
	scenarioLogger.Debug(msg, slog.Group("scenarioCtx",
		slog.String("scenarioRandID", string(scenarioInfo.ScenarioRandID)),
		slog.String("scenarioUri", sc.Uri),
		slog.String("scenarioName", sc.Name),
	))
}

func afterScenarioLogging(runLogger, suiteLogger, scenarioLogger *slog.Logger, sc *godog.Scenario, scenarioErr error, afterErr error, scenarioStart time.Time, scenarioInfo models_common.ScenarioInfo) {
	msg := "scenario success"
	elapsed := time.Since(scenarioStart)

	runGroup := slog.Group("scenarioCtx",
		slog.String("suiteName", string(scenarioInfo.SuiteName)),
		slog.String("scenarioRandID", string(scenarioInfo.ScenarioRandID)),
		slog.String("scenarioUri", sc.Uri),
		slog.String("scenarioName", sc.Name),
		slog.String("scenarioDuration", elapsed.String()),
	)

	suiteGroup := slog.Group("scenarioCtx",
		slog.String("scenarioRandID", string(scenarioInfo.ScenarioRandID)),
		slog.String("scenarioUri", sc.Uri),
		slog.String("scenarioName", sc.Name),
		slog.String("scenarioDuration", elapsed.String()),
	)

	scenarioGroup := slog.Group("scenarioCtx",
		slog.String("scenarioRandID", string(scenarioInfo.ScenarioRandID)),
		slog.String("scenarioUri", sc.Uri),
		slog.String("scenarioName", sc.Name),
		slog.String("scenarioDuration", elapsed.String()),
	)

	if scenarioErr != nil || afterErr != nil {
		msg = "scenario failed"
		slog.Error(msg,
			slog.Any(xtestlog.AttrErrorKey, fx.FlattenErrorsIfAny(scenarioErr, afterErr)),
			runGroup,
		)
		runLogger.Error(msg,
			slog.Any(xtestlog.AttrErrorKey, fx.FlattenErrorsIfAny(scenarioErr, afterErr)),
			runGroup,
		)
		suiteLogger.Error(msg,
			slog.Any(xtestlog.AttrErrorKey, fx.FlattenErrorsIfAny(scenarioErr, afterErr)),
			suiteGroup,
		)
		scenarioLogger.Error(msg,
			slog.Any(xtestlog.AttrErrorKey, fx.FlattenErrorsIfAny(scenarioErr, afterErr)),
			scenarioGroup,
		)
	} else {
		runLogger.Info(msg, runGroup)
		suiteLogger.Info(msg, suiteGroup)
		scenarioLogger.Info(msg, scenarioGroup)
	}
}
