//go:build web

package api

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"strings"
	"tinyrdm/backend/services"

	"github.com/gin-gonic/gin"
)

// maxRequestBodySize limits request body to 10MB to prevent memory exhaustion
const maxRequestBodySize = 10 << 20 // 10MB

// SetupRouter creates the Gin router with all API routes and static file serving
func SetupRouter(assets embed.FS) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Request body size limit
	r.Use(func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxRequestBodySize)
		c.Next()
	})

	// Security headers
	r.Use(SecurityHeaders())

	// CORS - validate origin for cross-origin requests
	r.Use(func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" {
			if isSameOrigin(c, origin) {
				c.Header("Access-Control-Allow-Origin", origin)
				c.Header("Access-Control-Allow-Credentials", "true")
				c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				c.Header("Access-Control-Allow-Headers", "Content-Type, X-Requested-With")
			} else {
				log.Printf("[cors] blocked origin=%s host=%s", origin, getRequestHost(c))
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
		}
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// CSRF protection for state-changing requests
	r.Use(csrfProtection())

	// Public routes (no auth required)
	registerAuthRoutes(r)
	r.GET("/api/version", func(c *gin.Context) {
		resp := services.Preferences().GetAppVersion()
		c.JSON(200, resp)
	})

	// WebSocket endpoint (auth checked via cookie + origin)
	r.GET("/ws", wsAuthCheck(), Hub().HandleWebSocket)

	// Protected API routes
	api := r.Group("/api")
	api.Use(AuthMiddleware())
	registerConnectionRoutes(api)
	registerBrowserRoutes(api)
	registerCLIRoutes(api)
	registerMonitorRoutes(api)
	registerPubsubRoutes(api)
	registerPreferencesRoutes(api)
	registerSystemRoutes(api)

	// Serve frontend static files from embedded assets
	distFS, err := fs.Sub(assets, "frontend/dist")
	if err == nil {
		fileServer := http.FileServer(http.FS(distFS))
		r.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path
			f, ferr := http.FS(distFS).Open(path)
			if ferr == nil {
				f.Close()
				fileServer.ServeHTTP(c.Writer, c.Request)
				return
			}
			c.FileFromFS("/", http.FS(distFS))
		})
	}

	return r
}

// getRequestHost returns the effective host, considering reverse proxy headers
func getRequestHost(c *gin.Context) string {
	if fwdHost := c.GetHeader("X-Forwarded-Host"); fwdHost != "" {
		return fwdHost
	}
	return c.Request.Host
}

// stripPort removes port from host string ("example.com:8088" -> "example.com")
func stripPort(host string) string {
	if idx := strings.LastIndex(host, ":"); idx >= 0 {
		// Make sure it's not part of IPv6 address
		if !strings.Contains(host, "]") || strings.LastIndex(host, "]") < idx {
			return host[:idx]
		}
	}
	return host
}

// extractOriginHost extracts hostname from Origin header value
func extractOriginHost(origin string) string {
	host := origin
	if idx := strings.Index(host, "://"); idx >= 0 {
		host = host[idx+3:]
	}
	host = strings.TrimRight(host, "/")
	return host
}

// isSameOrigin checks if the Origin header matches the request host.
// Compares hostnames only (ignoring port) to support reverse proxy scenarios
// where the external port differs from the internal port.
func isSameOrigin(c *gin.Context, origin string) bool {
	originHost := stripPort(extractOriginHost(origin))
	requestHost := stripPort(getRequestHost(c))
	return originHost == requestHost
}

// csrfProtection validates Origin/Referer for state-changing requests
func csrfProtection() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		if method == "GET" || method == "HEAD" || method == "OPTIONS" {
			c.Next()
			return
		}

		// Check Origin header first
		origin := c.GetHeader("Origin")
		if origin != "" {
			if !isSameOrigin(c, origin) {
				log.Printf("[csrf] blocked origin=%s host=%s", origin, getRequestHost(c))
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"success": false, "msg": "cross-origin request blocked"})
				return
			}
			c.Next()
			return
		}

		// Fallback: check Referer
		referer := c.GetHeader("Referer")
		if referer != "" {
			refererHost := extractOriginHost(referer)
			if slashIdx := strings.Index(refererHost, "/"); slashIdx >= 0 {
				refererHost = refererHost[:slashIdx]
			}
			requestHost := stripPort(getRequestHost(c))
			if stripPort(refererHost) != requestHost {
				log.Printf("[csrf] blocked referer=%s host=%s", referer, getRequestHost(c))
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"success": false, "msg": "cross-origin request blocked"})
				return
			}
		}

		c.Next()
	}
}

// wsAuthCheck validates auth and origin for WebSocket connections
func wsAuthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" {
			if !isSameOrigin(c, origin) {
				log.Printf("[ws] blocked origin=%s host=%s", origin, getRequestHost(c))
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
		}

		if !IsAuthEnabled() {
			c.Next()
			return
		}
		token, err := c.Cookie("rdm_token")
		if err != nil || !validateToken(token, getClientIP(c)) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Next()
	}
}
