package routers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/org/2112-space-lab/org/app-service/internal/api/handlers/errors"
	healthHandlers "github.com/org/2112-space-lab/org/app-service/internal/api/handlers/healthz"
	metricsHandlers "github.com/org/2112-space-lab/org/app-service/internal/api/handlers/metrics"
	"github.com/org/2112-space-lab/org/app-service/internal/config"
	"github.com/org/2112-space-lab/org/app-service/internal/services"
	logger "github.com/org/2112-space-lab/org/app-service/pkg/log"
	"github.com/prometheus/client_golang/prometheus"
)

// ProtectedRouter manages the public API router and its dependencies.
type ProtectedRouter struct {
	Echo             *echo.Echo
	Name             string
	ServiceComponent *services.ServiceComponent
}

// Init initializes the Echo instance for the router.
func (r *ProtectedRouter) Init() {
	r.Echo = echo.New()
	r.Echo.HideBanner = true

}

func InitProtectedAPIRouter(env *config.SEnv, services *services.ServiceComponent) *ProtectedRouter {
	logger.Debug("Initializing protected api router ...")
	protectedApiRouter := &ProtectedRouter{
		Name:             "public API",
		ServiceComponent: services,
	}
	protectedApiRouter.Init()
	protectedApiRouter.registerPrometheusMetrics()

	// next register all health check routes
	logger.Debug("Registering protected api health routes ...")
	protectedApiRouter.registerProtectedApiHealthCheckHandlers()

	logger.Debug("Registering metrics api protected routes ...")
	protectedApiRouter.registerMetricsAPIRoutes()

	// finally register default fallback error handlers
	// 404 is handled here as the last route
	logger.Debug("Registering protected api error handlers ...")
	protectedApiRouter.registerProtectedApiErrorHandlers()

	logger.Debug("Protected api registration complete.")
	return protectedApiRouter
}

func (r *ProtectedRouter) registerPrometheusMetrics() {
	appUptime := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "app_uptime_seconds",
			Help: "Application uptime in seconds.",
		},
	)
	appLiveness := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "app_liveness",
			Help: "Application liveness status (1 = alive, 0 = not alive).",
		},
	)

	prometheus.MustRegister(appUptime)
	prometheus.MustRegister(appLiveness)

	go func() {
		startTime := time.Now()
		appLiveness.Set(1)
		for {
			appUptime.Set(time.Since(startTime).Seconds())
			time.Sleep(1 * time.Second)
		}
	}()
}

func (r *ProtectedRouter) registerProtectedApiErrorHandlers() {
	r.Echo.HTTPErrorHandler = errors.AutomatedHttpErrorHandler()
	r.Echo.RouteNotFound("/*", errors.NotFound)
}

func (r *ProtectedRouter) registerProtectedApiHealthCheckHandlers() {
	health := r.Echo.Group("/health")
	health.GET("/alive", healthHandlers.Index)
	health.GET("/ready", healthHandlers.Ready)
}
func (r *ProtectedRouter) registerMetricsAPIRoutes() {
	metrics := r.Echo.Group("/metrics")
	metrics.GET("", metricsHandlers.GetMetrics)
}

// Start the Echo server
func (r *ProtectedRouter) Start(host string, port string) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		r.Echo.Logger.Info(fmt.Sprintf("Starting %s server on port: %s", r.Name, port))
		if err := r.Echo.Start(host + ":" + port); err != nil && err != http.ErrServerClosed {
			r.Echo.Logger.Fatalf("Server error: %v", err)
		}
	}()

	<-ctx.Done()
	r.shutdownServer()
}

func (r *ProtectedRouter) shutdownServer() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := r.Echo.Shutdown(ctx); err != nil {
		r.Echo.Logger.Fatalf("Failed to shutdown server: %v", err)
	}

	r.Echo.Logger.Info("Server shutdown complete.")
}
