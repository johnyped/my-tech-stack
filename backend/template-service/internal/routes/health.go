// Package routes registers HTTP routes on the Gin engine, grouped one file
// per domain — add internal/routes/example.go etc. alongside this one as you
// add domains, and call RegisterXxxRoutes from cmd/main.go's newRouter.
package routes

import (
	"github.com/gin-gonic/gin"

	"gitlab.awesome-poc-th.com/libero-engineering/TODO-group/TODO-service-name/internal/handler/health"
)

// RegisterHealthRoutes mounts the liveness, readiness, and ping probes.
func RegisterHealthRoutes(r *gin.Engine, h *health.Handler) {
	g := r.Group("/health")
	g.GET("/liveness", h.Liveness)
	g.GET("/readiness", h.Readiness)
	r.GET("/ping", func(c *gin.Context) { c.Status(200) })
}
