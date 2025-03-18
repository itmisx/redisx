package redisx

import (
	"context"
	"strings"

	"github.com/go-redis/redis/v8"
)

type prefixHook struct {
	Prefix string
}

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

// BeforeProcess 在执行 Redis 命令前，修改 key 以添加前缀
func (h *prefixHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	h.addPrefix(cmd)
	return ctx, nil
}
func (h *prefixHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	return nil
}

// BeforeProcessPipeline 处理管道命令
func (h *prefixHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	for _, cmd := range cmds {
		h.addPrefix(cmd)
	}
	return ctx, nil
}

// BeforeProcessPipeline 处理管道命令
func (h *prefixHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	return nil
}
