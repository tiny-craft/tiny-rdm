package services

import (
	"bytes"
	"encoding/json"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"tinyrdm/backend/storage"

	"github.com/google/uuid"
)

// google analytics service
type gaService struct {
	measurementID string
	secretKey     string
	clientID      string
}

type GaDataItem struct {
	ClientID string        `json:"client_id"`
	Events   []GaEventItem `json:"events"`
}

type GaEventItem struct {
	Name   string         `json:"name"`
	Params map[string]any `json:"params"`
}

var ga *gaService
var onceGA sync.Once

func GA() *gaService {
	if ga == nil {
		onceGA.Do(func() {
			// get or create an unique user id
			st := storage.NewLocalStore("device.txt")
			uidByte, err := st.Load()
			if err != nil {
				uidByte = []byte(strings.ReplaceAll(uuid.NewString(), "-", ""))
				st.Store(uidByte)
			}

			ga = &gaService{
				clientID: string(uidByte),
			}
		})
	}
	return ga
}

func (a *gaService) SetSecretKey(measurementID, secretKey string) {
	a.measurementID = measurementID
	a.secretKey = secretKey
}

func (a *gaService) isValid() bool {
	return len(a.measurementID) > 0 && len(a.secretKey) > 0
}

func (a *gaService) sendEvent(events ...GaEventItem) error {
	body, err := json.Marshal(GaDataItem{
		ClientID: a.clientID,
		Events:   events,
	})
	if err != nil {
		return err
	}

	//url := "https://www.google-analytics.com/debug/mp/collect"
	url := "https://www.google-analytics.com/mp/collect"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Add("measurement_id", a.measurementID)
	q.Add("api_secret", a.secretKey)
	req.URL.RawQuery = q.Encode()

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	//if dump, err := httputil.DumpResponse(response, true); err == nil {
	//	log.Println(string(dump))
	//}

	return nil
}

// Startup sends application startup event
func (a *gaService) Startup(version string) {
	if !a.isValid() {
		return
	}

	go a.sendEvent(GaEventItem{
		Name: "startup",
		Params: map[string]any{
			"os":      runtime.GOOS,
			"arch":    runtime.GOARCH,
			"version": version,
		},
	})
}
