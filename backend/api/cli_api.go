//go:build web

package api

import (
	"net/http"
	"tinyrdm/backend/services"
	"tinyrdm/backend/types"

	"github.com/gin-gonic/gin"
)

func registerCLIRoutes(rg *gin.RouterGroup) {
	g := rg.Group("/cli")

	g.POST("/start", func(c *gin.Context) {
		var req struct {
			Server string `json:"server"`
			DB     int    `json:"db"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, types.JSResp{Msg: "invalid request"})
			return
		}
		c.JSON(http.StatusOK, services.Cli().StartCli(req.Server, req.DB))
	})

	g.POST("/close", func(c *gin.Context) {
		var req struct {
			Server string `json:"server"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, types.JSResp{Msg: "invalid request"})
			return
		}
		c.JSON(http.StatusOK, services.Cli().CloseCli(req.Server))
	})

	// CLI input is handled via WebSocket - the frontend sends
	// {"event": "cmd:input:<server>", "data": "<command>"} over WS
}
