package routers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	apiaudittrail "github.com/Elbujito/2112/src/app-service/internal/api/handlers/audits"
	apicontext "github.com/Elbujito/2112/src/app-service/internal/api/handlers/context"
	"github.com/Elbujito/2112/src/app-service/internal/api/handlers/errors"
	healthHandlers "github.com/Elbujito/2112/src/app-service/internal/api/handlers/healthz"
	"github.com/Elbujito/2112/src/app-service/internal/api/handlers/satellites"
	"github.com/Elbujito/2112/src/app-service/internal/api/handlers/tiles"
	apiuser "github.com/Elbujito/2112/src/app-service/internal/api/handlers/users"
	"github.com/Elbujito/2112/src/app-service/internal/api/middlewares"
	"github.com/Elbujito/2112/src/app-service/internal/config"
	"github.com/Elbujito/2112/src/app-service/internal/services"
	logger "github.com/Elbujito/2112/src/app-service/pkg/log"
	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/labstack/echo/v4"
)

// PublicRouter manages the public API router and its dependencies.
type PublicRouter struct {
	Echo              *echo.Echo
	Name              string
	ServiceComponent  *services.ServiceComponent
	RouteTableMapping map[string]string
}

// NewPublicRouter creates and initializes a new PublicRouter instance.
func NewPublicRouter(env *config.SEnv, services *services.ServiceComponent) *PublicRouter {
	logger.Debug("Initializing public API router ...")
	// clerk.SetKey(env.EnvVars.Clerk.CLERK_API_KEY)
	clerk.SetKey("sk_test_GkqI0OhxlxMiywMZ2zgoNhGZ5H4RYymSdfDfdiTPBc")

	router := &PublicRouter{
		Name:             "public API",
		ServiceComponent: services,
	}
	router.setupEcho()
	router.registerRoutes()
	router.registerMiddlewares()

	logger.Debug("Public API router initialization complete.")
	return router
}

// setupEcho configures the Echo instance.
func (r *PublicRouter) setupEcho() {
	r.Echo = echo.New()
	r.Echo.HideBanner = true
}

// registerMiddlewares configures middleware for the Echo instance.
func (r *PublicRouter) registerMiddlewares() {
	logger.Debug("Registering middlewares...")
	middlewareList := []echo.MiddlewareFunc{
		middlewares.SlashesMiddleware(),
		middlewares.LoggerMiddleware(),
		middlewares.TimeoutMiddleware(),
		middlewares.ResponseHeadersMiddleware(),
		middlewares.ClerkMiddleware(),
		middlewares.LogNonGETRequestsMiddleware(r.RouteTableMapping, r.ServiceComponent.AuditTrailService),
	}

	if config.DevModeFlag {
		middlewareList = append(middlewareList, middlewares.BodyDumpMiddleware())
	}

	for _, middleware := range middlewareList {
		r.Echo.Use(middleware)
	}
	logger.Debug("Middlewares registered.")
}

// registerRoutes registers API routes for the Echo instance.
func (r *PublicRouter) registerRoutes() {
	logger.Debug("Registering routes...")
	r.registerHealthCheckRoutes()
	r.registerPublicAPIRoutes()
	r.registerErrorHandlers()
	logger.Debug("Routes registered.")
}

// registerHealthCheckRoutes registers health check routes.
func (r *PublicRouter) registerHealthCheckRoutes() {
	health := r.Echo.Group("/health")
	health.GET("/alive", healthHandlers.Index)
	health.GET("/ready", healthHandlers.Ready)
}

// registerErrorHandlers sets up global error handlers for the Echo instance.
func (r *PublicRouter) registerErrorHandlers() {
	r.Echo.HTTPErrorHandler = errors.AutomatedHttpErrorHandler()
	r.Echo.RouteNotFound("/*", errors.NotFound)
}

// registerPublicAPIRoutes registers public API routes.
func (r *PublicRouter) registerPublicAPIRoutes() {
	r.RouteTableMapping = middlewares.GenerateRouteTableMapping(r.Echo)

	// Handlers
	satelliteHandler := satellites.NewSatelliteHandler(r.ServiceComponent.SatelliteService)
	contextHandler := apicontext.NewContextHandler(r.ServiceComponent.ContextService)
	tileHandler := tiles.NewTileHandler(r.ServiceComponent.TileService)
	auditTrailHandler := apiaudittrail.NewAuditTrailHandler(r.ServiceComponent.AuditTrailService)
	userHandler := apiuser.NewUserHandler()

	// Satellite routes
	satellite := r.Echo.Group("/satellites")
	satellite.GET("/orbit", satelliteHandler.GetSatellitePositionsByNoradID)
	satellite.GET("/paginated", satelliteHandler.GetPaginatedSatellites)
	satellite.GET("/paginated/tles", satelliteHandler.GetPaginatedSatelliteInfo)

	// Tile routes
	tile := r.Echo.Group("/tiles")
	tile.GET("/all", tileHandler.GetAllTiles)
	tile.GET("/region", tileHandler.GetTilesInRegionHandler)
	tile.GET("/mappings", tileHandler.GetPaginatedSatelliteMappings)
	tile.PUT("/mappings/recompute/bynoradID", tileHandler.RecomputeMappingsByNoradID)
	tile.GET("/mappings/bynoradID", tileHandler.GetSatelliteMappingsByNoradID)

	// Context routes
	context := r.Echo.Group("/contexts")
	context.GET("/all", contextHandler.GetPaginatedContexts)
	context.POST("/", contextHandler.CreateContext)
	context.PUT("/:name", contextHandler.UpdateContext)
	context.GET("/:name", contextHandler.GetContextByName)
	context.DELETE("/:name", contextHandler.DeleteContextByName)
	context.PUT("/:name/activate", contextHandler.ActivateContext)
	context.PUT("/:name/deactivate", contextHandler.DeactivateContext)

	// Audit trail routes
	audit := r.Echo.Group("/audit-trails")
	audit.GET("/", auditTrailHandler.GetAuditTrails)

	// User routes
	user := r.Echo.Group("/users")
	user.GET("/", userHandler.GetUsers)
}

// Start runs the Echo server and handles graceful shutdown.
func (r *PublicRouter) Start(host, port string) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	serverAddress := fmt.Sprintf("%s:%s", host, port)
	go func() {
		logger.Infof("Starting %s server on %s", r.Name, serverAddress)
		if err := r.Echo.Start(serverAddress); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server error: %v", err)
		}
	}()

	<-ctx.Done()
	r.shutdownServer()
}

// shutdownServer handles the graceful shutdown of the Echo server.
func (r *PublicRouter) shutdownServer() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	logger.Info("Shutting down server...")
	if err := r.Echo.Shutdown(ctx); err != nil {
		logger.Fatal("Failed to shutdown server: %v", err)
	}
	logger.Info("Server shutdown complete.")
}
