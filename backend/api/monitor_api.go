//go:build web

package api

import (
	"net/http"
	"tinyrdm/backend/services"
	"tinyrdm/backend/types"

	"github.com/gin-gonic/gin"
)

func registerMonitorRoutes(rg *gin.RouterGroup) {
	g := rg.Group("/monitor")

	g.POST("/start", func(c *gin.Context) {
		var req struct {
			Server string `json:"server"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, types.JSResp{Msg: "invalid request"})
			return
		}
		c.JSON(http.StatusOK, services.Monitor().StartMonitor(req.Server))
	})

	g.POST("/stop", func(c *gin.Context) {
		var req struct {
			Server string `json:"server"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, types.JSResp{Msg: "invalid request"})
			return
		}
		c.JSON(http.StatusOK, services.Monitor().StopMonitor(req.Server))
	})

	g.POST("/export-log", func(c *gin.Context) {
		var req struct {
			Logs []string `json:"logs"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, types.JSResp{Msg: "invalid request"})
			return
		}
		c.JSON(http.StatusOK, services.Monitor().ExportLog(req.Logs))
	})
}
