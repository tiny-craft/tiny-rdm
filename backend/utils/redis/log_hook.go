package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"net"
	"strconv"
	"time"
)

type execCallback func(string, int64)

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

func appendArg(b []byte, v interface{}) []byte {
	switch v := v.(type) {
	case nil:
		return append(b, "<nil>"...)
	case string:
		return append(b, []byte(v)...)
	case []byte:
		return append(b, v...)
	case int:
		return strconv.AppendInt(b, int64(v), 10)
	case int8:
		return strconv.AppendInt(b, int64(v), 10)
	case int16:
		return strconv.AppendInt(b, int64(v), 10)
	case int32:
		return strconv.AppendInt(b, int64(v), 10)
	case int64:
		return strconv.AppendInt(b, v, 10)
	case uint:
		return strconv.AppendUint(b, uint64(v), 10)
	case uint8:
		return strconv.AppendUint(b, uint64(v), 10)
	case uint16:
		return strconv.AppendUint(b, uint64(v), 10)
	case uint32:
		return strconv.AppendUint(b, uint64(v), 10)
	case uint64:
		return strconv.AppendUint(b, v, 10)
	case float32:
		return strconv.AppendFloat(b, float64(v), 'f', -1, 64)
	case float64:
		return strconv.AppendFloat(b, v, 'f', -1, 64)
	case bool:
		if v {
			return append(b, "true"...)
		}
		return append(b, "false"...)
	case time.Time:
		return v.AppendFormat(b, time.RFC3339Nano)
	default:
		return append(b, fmt.Sprint(v)...)
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
		t := time.Now()
		err := next(ctx, cmd)
		if l.cmdExec != nil {
			b := make([]byte, 0, 64)
			for i, arg := range cmd.Args() {
				if i > 0 {
					b = append(b, ' ')
				}
				b = appendArg(b, arg)
			}
			l.cmdExec(string(b), time.Since(t).Milliseconds())
		}
		return err
	}
}

func (l *LogHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		t := time.Now()
		err := next(ctx, cmds)
		cost := time.Since(t).Milliseconds()
		for _, cmd := range cmds {
			log.Println("pipeline: ", cmd)
			if l.cmdExec != nil {
				b := make([]byte, 0, 64)
				for i, arg := range cmd.Args() {
					if i > 0 {
						b = append(b, ' ')
					}
					b = appendArg(b, arg)
				}
				l.cmdExec(string(b), cost)
			}
		}
		return err
	}
}
