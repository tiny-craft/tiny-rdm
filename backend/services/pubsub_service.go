package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"strconv"
	"sync"
	"time"
	"tinyrdm/backend/types"
)

type pubsubItem struct {
	client    *redis.Client
	pubsub    *redis.PubSub
	mutex     sync.Mutex
	closeCh   chan struct{}
	eventName string
}

type subMessage struct {
	Timestamp int64  `json:"timestamp"`
	Channel   string `json:"channel"`
	Message   string `json:"message"`
}

type pubsubService struct {
	ctx       context.Context
	ctxCancel context.CancelFunc
	mutex     sync.Mutex
	items     map[string]*pubsubItem
}

var pubsub *pubsubService
var oncePubsub sync.Once

func Pubsub() *pubsubService {
	if pubsub == nil {
		oncePubsub.Do(func() {
			pubsub = &pubsubService{
				items: map[string]*pubsubItem{},
			}
		})
	}
	return pubsub
}

func (p *pubsubService) getItem(server string) (*pubsubItem, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	item, ok := p.items[server]
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
		item = &pubsubItem{
			client: client,
		}
		p.items[server] = item
	}
	return item, nil
}

func (p *pubsubService) Start(ctx context.Context) {
	p.ctx, p.ctxCancel = context.WithCancel(ctx)
}

// Publish publish message to channel
func (p *pubsubService) Publish(server, channel, payload string) (resp types.JSResp) {
	rdb, err := Browser().getRedisClient(server, -1)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	var received int64
	received, err = rdb.client.Publish(p.ctx, channel, payload).Result()
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	resp.Success = true
	resp.Data = struct {
		Received int64 `json:"received"`
	}{
		Received: received,
	}
	return
}

// StartSubscribe start to subscribe a channel
func (p *pubsubService) StartSubscribe(server string) (resp types.JSResp) {
	item, err := p.getItem(server)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	item.closeCh = make(chan struct{})
	item.eventName = "sub:" + strconv.Itoa(int(time.Now().Unix()))
	item.pubsub = item.client.PSubscribe(p.ctx, "*")

	go p.processSubscribe(&item.mutex, item.pubsub.Channel(), item.closeCh, item.eventName)
	resp.Success = true
	resp.Data = struct {
		EventName string `json:"eventName"`
	}{
		EventName: item.eventName,
	}
	return
}

func (p *pubsubService) processSubscribe(mutex *sync.Mutex, ch <-chan *redis.Message, closeCh <-chan struct{}, eventName string) {
	lastEmitTime := time.Now().Add(-1 * time.Minute)
	cache := make([]subMessage, 0, 1000)
	for {
		select {
		case data := <-ch:
			go func() {
				timestamp := time.Now().UnixMilli()
				mutex.Lock()
				defer mutex.Unlock()
				cache = append(cache, subMessage{
					Timestamp: timestamp,
					Channel:   data.Channel,
					Message:   data.Payload,
				})
				if time.Now().Sub(lastEmitTime) > 300*time.Millisecond || len(cache) > 300 {
					runtime.EventsEmit(p.ctx, eventName, cache)
					cache = cache[:0:cap(cache)]
					lastEmitTime = time.Now()
				}
			}()

		case <-closeCh:
			// subscribe stopped
			return
		}
	}
}

// StopSubscribe stop subscribe by server name
func (p *pubsubService) StopSubscribe(server string) (resp types.JSResp) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	item, ok := p.items[server]
	if !ok || item.pubsub == nil {
		resp.Success = true
		return
	}

	//item.pubsub.Unsubscribe(p.ctx, "*")
	item.pubsub.Close()
	close(item.closeCh)
	delete(p.items, server)
	resp.Success = true
	return
}

// StopAll stop all subscribe
func (p *pubsubService) StopAll() {
	if p.ctxCancel != nil {
		p.ctxCancel()
	}

	for server := range p.items {
		p.StopSubscribe(server)
	}
}
