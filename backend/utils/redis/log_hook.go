package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"net"
)

type execCallback func(string)

type LogHook struct {
	name    string
	cmdExec execCallback
}

func NewHook(name string, cmdExec execCallback) *LogHook {
	return &LogHook{
		name:    name,
		cmdExec: cmdExec,
	}
}

func (l *LogHook) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		return next(ctx, network, addr)
	}
}
func (l *LogHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		log.Println(cmd)
		if l.cmdExec != nil {
			l.cmdExec(cmd.String())
		}
		return next(ctx, cmd)
	}
}
func (l *LogHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		for _, cmd := range cmds {
			log.Println("pipeline: ", cmd)
			if l.cmdExec != nil {
				l.cmdExec(cmd.String())
			}
		}
		return next(ctx, cmds)
	}
}
