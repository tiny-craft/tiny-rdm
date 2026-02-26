//go:build web

package api

import (
	"net/http"
	"tinyrdm/backend/services"
	"tinyrdm/backend/types"

	"github.com/gin-gonic/gin"
)

func registerBrowserRoutes(rg *gin.RouterGroup) {
	g := rg.Group("/browser")

	g.POST("/open-connection", func(c *gin.Context) {
		var req struct {
			Name string `json:"name"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, types.JSResp{Msg: "invalid request"})
			return
		}
		c.JSON(http.StatusOK, services.Browser().OpenConnection(req.Name))
	})

	g.POST("/close-connection", func(c *gin.Context) {
		var req struct {
			Name string `json:"name"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, types.JSResp{Msg: "invalid request"})
			return
		}
		c.JSON(http.StatusOK, services.Browser().CloseConnection(req.Name))
	})

	g.POST("/open-database", func(c *gin.Context) {
		var req struct {
			Server string `json:"server"`
			DB     int    `json:"db"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, types.JSResp{Msg: "invalid request"})
			return
		}
		c.JSON(http.StatusOK, services.Browser().OpenDatabase(req.Server, req.DB))
	})

	g.POST("/server-info", func(c *gin.Context) {
		var req struct {
			Name string `json:"name"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, types.JSResp{Msg: "invalid request"})
			return
		}
		c.JSON(http.StatusOK, services.Browser().ServerInfo(req.Name))
	})

	g.POST("/load-next-keys", func(c *gin.Context) {
		var req struct {
			Server     string `json:"server"`
			DB         int    `json:"db"`
			Match      string `json:"match"`
			KeyType    string `json:"keyType"`
			ExactMatch bool   `json:"exactMatch"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, types.JSResp{Msg: "invalid request"})
			return
		}
		c.JSON(http.StatusOK, services.Browser().LoadNextKeys(req.Server, req.DB, req.Match, req.KeyType, req.ExactMatch))
	})

	g.POST("/load-next-all-keys", func(c *gin.Context) {
		var req struct {
			Server     string `json:"server"`
			DB         int    `json:"db"`
			Match      string `json:"match"`
			KeyType    string `json:"keyType"`
			ExactMatch bool   `json:"exactMatch"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, types.JSResp{Msg: "invalid request"})
			return
		}
		c.JSON(http.StatusOK, services.Browser().LoadNextAllKeys(req.Server, req.DB, req.Match, req.KeyType, req.ExactMatch))
	})

	g.POST("/load-all-keys", func(c *gin.Context) {
		var req struct {
			Server     string `json:"server"`
			DB         int    `json:"db"`
			Match      string `json:"match"`
			KeyType    string `json:"keyType"`
			ExactMatch bool   `json:"exactMatch"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, types.JSResp{Msg: "invalid request"})
			return
		}
		c.JSON(http.StatusOK, services.Browser().LoadAllKeys(req.Server, req.DB, req.Match, req.KeyType, req.ExactMatch))
	})

	g.POST("/get-key-type", func(c *gin.Context) {
		var param types.KeySummaryParam
		if err := c.ShouldBindJSON(&param); err != nil {
			c.JSON(http.StatusBadRequest, types.JSResp{Msg: "invalid request"})
			return
		}
		c.JSON(http.StatusOK, services.Browser().GetKeyType(param))
	})

	g.POST("/get-key-summary", func(c *gin.Context) {
		var param types.KeySummaryParam
		if err := c.ShouldBindJSON(&param); err != nil {
			c.JSON(http.StatusBadRequest, types.JSResp{Msg: "invalid request"})
			return
		}
		c.JSON(http.StatusOK, services.Browser().GetKeySummary(param))
	})

	g.POST("/get-key-detail", func(c *gin.Context) {
		var param types.KeyDetailParam
		if err := c.ShouldBindJSON(&param); err != nil {
			c.JSON(http.StatusBadRequest, types.JSResp{Msg: "invalid request"})
			return
		}
		c.JSON(http.StatusOK, services.Browser().GetKeyDetail(param))
	})

	g.POST("/convert-value", func(c *gin.Context) {
		var req struct {
			Value  any    `json:"value"`
			Decode string `json:"decode"`
			Format string `json:"format"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, types.JSResp{Msg: "invalid request"})
			return
		}
		c.JSON(http.StatusOK, services.Browser().ConvertValue(req.Value, req.Decode, req.Format))
	})

	g.POST("/set-key-value", func(c *gin.Context) {
		var param types.SetKeyParam
		if err := c.ShouldBindJSON(&param); err != nil {
			c.JSON(http.StatusBadRequest, types.JSResp{Msg: "invalid request"})
			return
		}
		c.JSON(http.StatusOK, services.Browser().SetKeyValue(param))
	})

	g.POST("/get-hash-value", func(c *gin.Context) {
		var param types.GetHashParam
		if err := c.ShouldBindJSON(&param); err != nil {
			c.JSON(http.StatusBadRequest, types.JSResp{Msg: "invalid request"})
			return
		}
		c.JSON(http.StatusOK, services.Browser().GetHashValue(param))
	})

	g.POST("/set-hash-value", func(c *gin.Context) {
		var param types.SetHashParam
		if err := c.ShouldBindJSON(&param); err != nil {
			c.JSON(http.StatusBadRequest, types.JSResp{Msg: "invalid request"})
			return
		}
		c.JSON(http.StatusOK, services.Browser().SetHashValue(param))
	})

	g.POST("/add-hash-field", func(c *gin.Context) {
		var req struct {
			Server     string `json:"server"`
			DB         int    `json:"db"`
			Key        any    `json:"key"`
			Action     int    `json:"action"`
			FieldItems []any  `json:"fieldItems"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, types.JSResp{Msg: "invalid request"})
			return
		}
		c.JSON(http.StatusOK, services.Browser().AddHashField(req.Server, req.DB, req.Key, req.Action, req.FieldItems))
	})

	g.POST("/add-list-item", func(c *gin.Context) {
		var req struct {
			Server string `json:"server"`
			DB     int    `json:"db"`
			Key    any    `json:"key"`
			Action int    `json:"action"`
			Items  []any  `json:"items"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, types.JSResp{Msg: "invalid request"})
			return
		}
		c.JSON(http.StatusOK, services.Browser().AddListItem(req.Server, req.DB, req.Key, req.Action, req.Items))
	})

	g.POST("/set-list-item", func(c *gin.Context) {
		var param types.SetListParam
		if err := c.ShouldBindJSON(&param); err != nil {
			c.JSON(http.StatusBadRequest, types.JSResp{Msg: "invalid request"})
			return
		}
		c.JSON(http.StatusOK, services.Browser().SetListItem(param))
	})

	g.POST("/set-set-item", func(c *gin.Context) {
		var req struct {
			Server  string `json:"server"`
			DB      int    `json:"db"`
			Key     any    `json:"key"`
			Remove  bool   `json:"remove"`
			Members []any  `json:"members"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, types.JSResp{Msg: "invalid request"})
			return
		}
		c.JSON(http.StatusOK, services.Browser().SetSetItem(req.Server, req.DB, req.Key, req.Remove, req.Members))
	})

	g.POST("/update-set-item", func(c *gin.Context) {
		var param types.SetSetParam
		if err := c.ShouldBindJSON(&param); err != nil {
			c.JSON(http.StatusBadRequest, types.JSResp{Msg: "invalid request"})
			return
		}
		c.JSON(http.StatusOK, services.Browser().UpdateSetItem(param))
	})

	g.POST("/update-zset-value", func(c *gin.Context) {
		var param types.SetZSetParam
		if err := c.ShouldBindJSON(&param); err != nil {
			c.JSON(http.StatusBadRequest, types.JSResp{Msg: "invalid request"})
			return
		}
		c.JSON(http.StatusOK, services.Browser().UpdateZSetValue(param))
	})

	g.POST("/add-zset-value", func(c *gin.Context) {
		var req struct {
			Server     string             `json:"server"`
			DB         int                `json:"db"`
			Key        any                `json:"key"`
			Action     int                `json:"action"`
			ValueScore map[string]float64 `json:"valueScore"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, types.JSResp{Msg: "invalid request"})
			return
		}
		c.JSON(http.StatusOK, services.Browser().AddZSetValue(req.Server, req.DB, req.Key, req.Action, req.ValueScore))
	})

	g.POST("/add-stream-value", func(c *gin.Context) {
		var req struct {
			Server     string `json:"server"`
			DB         int    `json:"db"`
			Key        any    `json:"key"`
			ID         string `json:"id"`
			FieldItems []any  `json:"fieldItems"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, types.JSResp{Msg: "invalid request"})
			return
		}
		c.JSON(http.StatusOK, services.Browser().AddStreamValue(req.Server, req.DB, req.Key, req.ID, req.FieldItems))
	})

	g.POST("/remove-stream-values", func(c *gin.Context) {
		var req struct {
			Server string   `json:"server"`
			DB     int      `json:"db"`
			Key    any      `json:"key"`
			IDs    []string `json:"ids"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, types.JSResp{Msg: "invalid request"})
			return
		}
		c.JSON(http.StatusOK, services.Browser().RemoveStreamValues(req.Server, req.DB, req.Key, req.IDs))
	})

	g.POST("/set-key-ttl", func(c *gin.Context) {
		var req struct {
			Server string `json:"server"`
			DB     int    `json:"db"`
			Key    any    `json:"key"`
			TTL    int64  `json:"ttl"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, types.JSResp{Msg: "invalid request"})
			return
		}
		c.JSON(http.StatusOK, services.Browser().SetKeyTTL(req.Server, req.DB, req.Key, req.TTL))
	})

	g.POST("/batch-set-ttl", func(c *gin.Context) {
		var req struct {
			Server   string `json:"server"`
			DB       int    `json:"db"`
			Keys     []any  `json:"keys"`
			TTL      int64  `json:"ttl"`
			SerialNo string `json:"serialNo"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, types.JSResp{Msg: "invalid request"})
			return
		}
		c.JSON(http.StatusOK, services.Browser().BatchSetTTL(req.Server, req.DB, req.Keys, req.TTL, req.SerialNo))
	})

	g.POST("/delete-key", func(c *gin.Context) {
		var req struct {
			Server string `json:"server"`
			DB     int    `json:"db"`
			Key    any    `json:"key"`
			Async  bool   `json:"async"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, types.JSResp{Msg: "invalid request"})
			return
		}
		c.JSON(http.StatusOK, services.Browser().DeleteKey(req.Server, req.DB, req.Key, req.Async))
	})

	g.POST("/delete-keys", func(c *gin.Context) {
		var req struct {
			Server   string `json:"server"`
			DB       int    `json:"db"`
			Keys     []any  `json:"keys"`
			SerialNo string `json:"serialNo"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, types.JSResp{Msg: "invalid request"})
			return
		}
		c.JSON(http.StatusOK, services.Browser().DeleteKeys(req.Server, req.DB, req.Keys, req.SerialNo))
	})

	g.POST("/delete-keys-by-pattern", func(c *gin.Context) {
		var req struct {
			Server  string `json:"server"`
			DB      int    `json:"db"`
			Pattern string `json:"pattern"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, types.JSResp{Msg: "invalid request"})
			return
		}
		c.JSON(http.StatusOK, services.Browser().DeleteKeysByPattern(req.Server, req.DB, req.Pattern))
	})

	g.POST("/rename-key", func(c *gin.Context) {
		var req struct {
			Server string `json:"server"`
			DB     int    `json:"db"`
			Key    string `json:"key"`
			NewKey string `json:"newKey"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, types.JSResp{Msg: "invalid request"})
			return
		}
		c.JSON(http.StatusOK, services.Browser().RenameKey(req.Server, req.DB, req.Key, req.NewKey))
	})

	g.POST("/export-key", func(c *gin.Context) {
		var req struct {
			Server        string `json:"server"`
			DB            int    `json:"db"`
			Keys          []any  `json:"keys"`
			Path          string `json:"path"`
			IncludeExpire bool   `json:"includeExpire"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, types.JSResp{Msg: "invalid request"})
			return
		}
		c.JSON(http.StatusOK, services.Browser().ExportKey(req.Server, req.DB, req.Keys, req.Path, req.IncludeExpire))
	})

	g.POST("/import-csv", func(c *gin.Context) {
		var req struct {
			Server   string `json:"server"`
			DB       int    `json:"db"`
			Path     string `json:"path"`
			Conflict int    `json:"conflict"`
			TTL      int64  `json:"ttl"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, types.JSResp{Msg: "invalid request"})
			return
		}
		c.JSON(http.StatusOK, services.Browser().ImportCSV(req.Server, req.DB, req.Path, req.Conflict, req.TTL))
	})

	g.POST("/flush-db", func(c *gin.Context) {
		var req struct {
			Server string `json:"server"`
			DB     int    `json:"db"`
			Async  bool   `json:"async"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, types.JSResp{Msg: "invalid request"})
			return
		}
		c.JSON(http.StatusOK, services.Browser().FlushDB(req.Server, req.DB, req.Async))
	})

	g.POST("/get-slow-logs", func(c *gin.Context) {
		var req struct {
			Server string `json:"server"`
			Num    int64  `json:"num"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, types.JSResp{Msg: "invalid request"})
			return
		}
		c.JSON(http.StatusOK, services.Browser().GetSlowLogs(req.Server, req.Num))
	})

	g.POST("/get-client-list", func(c *gin.Context) {
		var req struct {
			Server string `json:"server"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, types.JSResp{Msg: "invalid request"})
			return
		}
		c.JSON(http.StatusOK, services.Browser().GetClientList(req.Server))
	})

	g.POST("/get-cmd-history", func(c *gin.Context) {
		var req struct {
			PageNo   int `json:"pageNo"`
			PageSize int `json:"pageSize"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, types.JSResp{Msg: "invalid request"})
			return
		}
		if req.PageSize <= 0 {
			req.PageSize = 50
		}
		c.JSON(http.StatusOK, services.Browser().GetCmdHistory(req.PageNo, req.PageSize))
	})

	g.POST("/clean-cmd-history", func(c *gin.Context) {
		c.JSON(http.StatusOK, services.Browser().CleanCmdHistory())
	})
}
