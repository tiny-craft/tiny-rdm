//go:build web

package services

import (
	"bytes"
	"io"
	"os"
	"path"
	"time"
	"tinyrdm/backend/types"

	"github.com/klauspost/compress/zip"
	"github.com/vrischmann/userdir"
)

// ExportConnectionsToBytes exports connections as zip bytes for web download
func (c *connectionService) ExportConnectionsToBytes() ([]byte, string, error) {
	const connectionFilename = "connections.yaml"
	filename := "connections_" + time.Now().Format("20060102150405") + ".zip"

	inputFile, err := os.Open(path.Join(userdir.GetConfigHome(), "TinyRDM", connectionFilename))
	if err != nil {
		return nil, "", err
	}
	defer inputFile.Close()

	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	headerWriter, err := zipWriter.CreateHeader(&zip.FileHeader{
		Name:   connectionFilename,
		Method: zip.Deflate,
	})
	if err != nil {
		return nil, "", err
	}

	if _, err = io.Copy(headerWriter, inputFile); err != nil {
		return nil, "", err
	}

	if err = zipWriter.Close(); err != nil {
		return nil, "", err
	}

	return buf.Bytes(), filename, nil
}

// ImportConnectionsFromBytes imports connections from uploaded zip bytes
func (c *connectionService) ImportConnectionsFromBytes(data []byte) (resp types.JSResp) {
	const connectionFilename = "connections.yaml"

	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		resp.Msg = "invalid zip file"
		return
	}

	var file *zip.File
	for _, f := range reader.File {
		if f.Name == connectionFilename {
			file = f
			break
		}
	}

	if file == nil {
		resp.Msg = "connections.yaml not found in zip"
		return
	}

	zippedFile, err := file.Open()
	if err != nil {
		resp.Msg = "failed to read zip content"
		return
	}
	defer zippedFile.Close()

	outputFile, err := os.Create(path.Join(userdir.GetConfigHome(), "TinyRDM", connectionFilename))
	if err != nil {
		resp.Msg = "failed to save connections"
		return
	}
	defer outputFile.Close()

	if _, err = io.Copy(outputFile, zippedFile); err != nil {
		resp.Msg = "failed to write connections"
		return
	}

	resp.Success = true
	return
}
