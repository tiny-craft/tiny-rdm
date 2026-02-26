//go:build web

package api

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"tinyrdm/backend/services"
	"tinyrdm/backend/types"

	"github.com/gin-gonic/gin"
)

// safeTempPath validates that a path is within the OS temp directory.
// Prevents directory traversal attacks.
func safeTempPath(reqPath string) (string, error) {
	tmpDir := os.TempDir()
	cleaned := filepath.Clean(reqPath)
	abs, err := filepath.Abs(cleaned)
	if err != nil {
		return "", fmt.Errorf("invalid path")
	}
	// Ensure the resolved path is within tmpDir
	if !strings.HasPrefix(abs, filepath.Clean(tmpDir)+string(os.PathSeparator)) && abs != filepath.Clean(tmpDir) {
		return "", fmt.Errorf("access denied")
	}
	return abs, nil
}

// sanitizeFilename removes path separators and dangerous characters from filename
func sanitizeFilename(name string) string {
	// Take only the base name to strip any directory components
	name = filepath.Base(name)
	// Remove any remaining path separators (extra safety)
	name = strings.ReplaceAll(name, "..", "")
	name = strings.ReplaceAll(name, "/", "")
	name = strings.ReplaceAll(name, "\\", "")
	if name == "" || name == "." {
		name = "upload"
	}
	return name
}

func registerSystemRoutes(rg *gin.RouterGroup) {
	g := rg.Group("/system")

	g.GET("/info", func(c *gin.Context) {
		c.JSON(http.StatusOK, services.System().Info())
	})

	// Web replacement for native file dialog - select file
	g.POST("/select-file", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, types.JSResp{Msg: "invalid file upload"})
			return
		}

		// Sanitize filename to prevent path traversal
		safeName := sanitizeFilename(file.Filename)
		tmpDir := os.TempDir()
		dst := filepath.Join(tmpDir, safeName)

		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.JSON(http.StatusInternalServerError, types.JSResp{Msg: "failed to save file"})
			return
		}

		c.JSON(http.StatusOK, types.JSResp{
			Success: true,
			Data: map[string]any{
				"path": dst,
			},
		})
	})

	// Web replacement for native file dialog - download file
	g.GET("/download", func(c *gin.Context) {
		reqPath := c.Query("path")
		if reqPath == "" {
			c.JSON(http.StatusBadRequest, types.JSResp{Msg: "path is required"})
			return
		}

		// Validate path is within temp directory only
		safePath, err := safeTempPath(reqPath)
		if err != nil {
			c.JSON(http.StatusForbidden, types.JSResp{Msg: "access denied"})
			return
		}

		file, err := os.Open(safePath)
		if err != nil {
			c.JSON(http.StatusNotFound, types.JSResp{Msg: "file not found"})
			return
		}
		defer file.Close()

		stat, err := file.Stat()
		if err != nil {
			c.JSON(http.StatusInternalServerError, types.JSResp{Msg: "failed to read file"})
			return
		}

		c.Header("Content-Disposition", "attachment; filename="+filepath.Base(safePath))
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Length", fmt.Sprintf("%d", stat.Size()))
		io.Copy(c.Writer, file)
	})
}
