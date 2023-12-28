package services

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"os"
	"strconv"
	"sync"
	"time"
	"tinyrdm/backend/types"
)

type monitorItem struct {
	client    *redis.Client
	cmd       *redis.MonitorCmd
	ch        chan string
	closeCh   chan struct{}
	eventName string
}

type monitorService struct {
	ctx       context.Context
	ctxCancel context.CancelFunc
	mutex     sync.Mutex
	items     map[string]*monitorItem
}

var monitor *monitorService
var onceMonitor sync.Once

func Monitor() *monitorService {
	if monitor == nil {
		onceMonitor.Do(func() {
			monitor = &monitorService{
				items: map[string]*monitorItem{},
			}
		})
	}
	return monitor
}

func (c *monitorService) getItem(server string) (*monitorItem, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	item, ok := c.items[server]
	if !ok {
		var err error
		conf := Connection().getConnection(server)
		if conf == nil {
			return nil, fmt.Errorf("no connection profile named: %s", server)
		}
		var uniClient redis.UniversalClient
		if uniClient, err = Connection().createRedisClient(conf.ConnectionConfig); err != nil {
			return nil, err
		}
		var client *redis.Client
		if client, ok = uniClient.(*redis.Client); !ok {
			return nil, errors.New("create redis client fail")
		}
		item = &monitorItem{
			client: client,
		}
		c.items[server] = item
	}
	return item, nil
}

func (c *monitorService) Start(ctx context.Context) {
	c.ctx, c.ctxCancel = context.WithCancel(ctx)
}

// StartMonitor start a monitor by server name
func (c *monitorService) StartMonitor(server string) (resp types.JSResp) {
	item, err := c.getItem(server)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	item.ch = make(chan string)
	item.closeCh = make(chan struct{})
	item.eventName = "monitor:" + strconv.Itoa(int(time.Now().Unix()))
	item.cmd = item.client.Monitor(c.ctx, item.ch)
	item.cmd.Start()

	go c.processMonitor(item.ch, item.closeCh, item.eventName)
	resp.Success = true
	resp.Data = struct {
		EventName string `json:"eventName"`
	}{
		EventName: item.eventName,
	}
	return
}

func (c *monitorService) processMonitor(ch <-chan string, closeCh <-chan struct{}, eventName string) {
	for {
		select {
		case data := <-ch:
			if data != "OK" {
				runtime.EventsEmit(c.ctx, eventName, data)
			}

		case <-closeCh:
			// monitor stopped
			return
		}
	}
}

// StopMonitor stop monitor by server name
func (c *monitorService) StopMonitor(server string) (resp types.JSResp) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	item, ok := c.items[server]
	if !ok || item.cmd == nil {
		resp.Success = true
		return
	}

	item.cmd.Stop()
	//close(item.ch)
	close(item.closeCh)
	delete(c.items, server)
	resp.Success = true
	return
}

// StopAll stop all monitor
func (c *monitorService) StopAll() {
	if c.ctxCancel != nil {
		c.ctxCancel()
	}

	for server := range c.items {
		c.StopMonitor(server)
	}
}

func (c *monitorService) ExportLog(logs []string) (resp types.JSResp) {
	filepath, err := runtime.SaveFileDialog(c.ctx, runtime.SaveDialogOptions{
		ShowHiddenFiles: false,
		DefaultFilename: fmt.Sprintf("monitor_log_%s.txt", time.Now().Format("20060102150405")),
		Filters: []runtime.FileFilter{
			{Pattern: "*.txt"},
		},
	})
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	file, err := os.Create(filepath)
	if err != nil {
		resp.Msg = err.Error()
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range logs {
		_, _ = writer.WriteString(line + "\n")
	}
	writer.Flush()

	resp.Success = true
	return
}
