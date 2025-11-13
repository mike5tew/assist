package api

import (
	"esp-organizer/internal/InfoFlow/InfoIn/api"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterDiagnosticsRoutes registers diagnostic endpoints for debugging and monitoring
func RegisterDiagnosticsRoutes(router *gin.Engine) {
	diagnostics := router.Group("/api/diagnostics")
	{
		// Use handlers from the consolidated common_handlers.go file
		diagnostics.GET("/config", wrapHandler(api.DiagnosticsConfigHandler))
		diagnostics.GET("/document", wrapHandler(api.DiagnosticsDocumentHandler))
		diagnostics.GET("/collections", wrapHandler(api.DiagnosticsCollectionsHandler))
	}
}

// wrapHandler converts http.HandlerFunc to gin.HandlerFunc
func wrapHandler(handler http.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler(c.Writer, c.Request)
	}
}
