package redisx

import (
	"context"
	"net"
	"strings"

	"github.com/redis/go-redis/v9"
)

// 支持多个 key 的命令
var multiKeyCmds = map[string]bool{
	"DEL":    true,
	"MGET":   true,
	"EXISTS": true,
	"UNLINK": true,
	"TOUCH":  true,
}

// 多 key + value 的命令
var multiKeyValueCmds = map[string]bool{
	"MSET": true,
}

// 统一处理 key 前缀
func (h *prefixHook) addPrefix(cmd redis.Cmder) {
	if len(cmd.Args()) < 2 {
		return
	}

	cmdName := strings.ToUpper(cmd.Args()[0].(string))

	// 处理多个 key 的命令
	if multiKeyCmds[cmdName] {
		for i := 1; i < len(cmd.Args()); i++ {
			if key, ok := cmd.Args()[i].(string); ok {
				cmd.Args()[i] = h.Prefix + key
			}
		}
		return
	}

	// 处理 key-value 组合的命令（MSET key1 val1 key2 val2）
	if multiKeyValueCmds[cmdName] {
		for i := 1; i < len(cmd.Args()); i += 2 {
			if key, ok := cmd.Args()[i].(string); ok {
				cmd.Args()[i] = h.Prefix + key
			}
		}
		return
	}

	// 默认 key 在第二个参数
	if key, ok := cmd.Args()[1].(string); ok {
		cmd.Args()[1] = h.Prefix + key
	}
}

type prefixHook struct {
	Prefix string
}

func (h *prefixHook) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		return next(ctx, network, addr)
	}
}

func (h *prefixHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		// 修改命令参数
		if len(cmd.Args()) > 1 {
			if key, ok := cmd.Args()[1].(string); ok {
				cmd.Args()[1] = h.Prefix + key
			}
		}
		return next(ctx, cmd) // 继续执行命令
	}
}

func (h *prefixHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		// 修改所有 Pipeline 里的 Key
		for _, cmd := range cmds {
			if len(cmd.Args()) > 1 {
				if key, ok := cmd.Args()[1].(string); ok {
					cmd.Args()[1] = h.Prefix + key
				}
			}
		}
		return next(ctx, cmds)
	}
}
