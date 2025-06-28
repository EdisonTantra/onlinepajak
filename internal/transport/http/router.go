package http

import (
	"net/http"

	"github.com/EdisonTantra/lemonPajak/internal/core/cons"
	lemonPort "github.com/EdisonTantra/lemonPajak/internal/core/port"
	"github.com/EdisonTantra/lemonPajak/pkg/lib/logat"
	lemonTracer "github.com/EdisonTantra/lemonPajak/pkg/lib/tracer"
	"github.com/gin-gonic/gin"
)

// Router represents http routing for managing apps services.
type Router struct {
	logger      logat.AppsLogger
	appHandler  lemonPort.AppsHandler
	userHandler lemonPort.AppsHandler
}

// NewRouter creates and returns a new instance of the App Router.
// It takes services as parameters and returns a pointer to the Router struct.
func NewRouter(
	logger logat.AppsLogger,
	appHandler lemonPort.AppsHandler,
	userHandler lemonPort.AppsHandler,
) *Router {
	return &Router{
		logger:      logger,
		appHandler:  appHandler,
		userHandler: userHandler,
	}
}

// Handlers returns list of http handler that mapped into an endpoint.
// It takes no parameters and returns http.Handler
func (r *Router) Handlers(serviceName string) http.Handler {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// register middleware
	router.Use(gin.Recovery())
	router.Use(r.MiddlewareLogging())
	router.GET("/health", r.Health())

	r.appHandler.Mount(router)
	r.userHandler.Mount(router)

	return router
}

// Health handle health check of service running
// with stable database connection
func (r *Router) Health() gin.HandlerFunc {
	return func(gCtx *gin.Context) {
		ctx := gCtx.Request.Context()
		tr := lemonTracer.StartTrace(ctx, cons.EventLogNameHealth)
		defer tr.Finish()

		gCtx.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	}
}
