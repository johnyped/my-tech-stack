// Package health implements the liveness and readiness probes every service
// exposes — see document/ARCHITECTURE.md "Health checks". Responses are
// use a minimal JSON contract consumed by Docker and Kubernetes probes.
package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	statusAlive = "ALIVE"
	statusReady = "READY"
	// statusNotReady = "NOT_READY" // use this once Readiness checks real dependencies
)

// Response is the shared health-check contract.
type Response struct {
	Status    string `json:"status"`
	ServiceID string `json:"service_id"`
}

// Handler serves GET /health/liveness and GET /health/readiness.
type Handler struct {
	serviceID string
	// Add dependency handles here (cache, db) once this service has any —
	// see Readiness below for where they get pinged.
}

// NewHandler constructs the health Handler.
func NewHandler(serviceID string) *Handler {
	return &Handler{serviceID: serviceID}
}

// Liveness reports whether the process itself is up. Never checks
// dependencies — that's Readiness's job. Used by Docker HEALTHCHECK.
func (h *Handler) Liveness(c *gin.Context) {
	c.JSON(http.StatusOK, Response{Status: statusAlive, ServiceID: h.serviceID})
}

// Readiness reports whether this instance can serve traffic — i.e. its
// configured dependencies (cache, db, kafka) are reachable. TODO: once this
// service has a dependency, ping it here with a short timeout (2s is the
// convention, see ekyc-service's health handler) and return NOT_READY /
// 503 on failure instead of always returning READY.
func (h *Handler) Readiness(c *gin.Context) {
	c.JSON(http.StatusOK, Response{Status: statusReady, ServiceID: h.serviceID})
}
