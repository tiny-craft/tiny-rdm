package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"net"
)

type LogHook struct {
	name string
}

func NewHook(name string) LogHook {
	return LogHook{
		name: name,
	}
}

func (LogHook) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		return next(ctx, network, addr)
	}
}
func (LogHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		log.Println(cmd.String())
		return next(ctx, cmd)
	}
}
func (LogHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		for _, cmd := range cmds {
			log.Println(cmd.String())
		}
		return next(ctx, cmds)
	}
}
