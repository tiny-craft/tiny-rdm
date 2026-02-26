//go:build web

package api

import (
	"net/http"
	"tinyrdm/backend/services"
	"tinyrdm/backend/types"

	"github.com/gin-gonic/gin"
)

func registerPubsubRoutes(rg *gin.RouterGroup) {
	g := rg.Group("/pubsub")

	g.POST("/publish", func(c *gin.Context) {
		var req struct {
			Server  string `json:"server"`
			Channel string `json:"channel"`
			Payload string `json:"payload"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, types.JSResp{Msg: "invalid request"})
			return
		}
		c.JSON(http.StatusOK, services.Pubsub().Publish(req.Server, req.Channel, req.Payload))
	})

	g.POST("/subscribe", func(c *gin.Context) {
		var req struct {
			Server string `json:"server"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, types.JSResp{Msg: "invalid request"})
			return
		}
		c.JSON(http.StatusOK, services.Pubsub().StartSubscribe(req.Server))
	})

	g.POST("/unsubscribe", func(c *gin.Context) {
		var req struct {
			Server string `json:"server"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, types.JSResp{Msg: "invalid request"})
			return
		}
		c.JSON(http.StatusOK, services.Pubsub().StopSubscribe(req.Server))
	})
}
