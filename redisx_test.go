package redisx_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/itmisx/redisx"
)

func TestRedisx(*testing.T) {
	// 实例变量
	var cli redisx.Client
	// 配置
	config := redisx.Config{
		Host:     "127.0.0.1",
		Port:     "6379",
		Protocol: "tcp",
	}
	// 新建实例
	cli = redisx.New(config)

	// 命令测试

	// SetEx
	{
		fmt.Println("----Set----")
		r, err := cli.SetEX(context.Background(), "key", "123", time.Second*10).Result()
		fmt.Println("result", r)
		fmt.Println("err", err)
	}

	// Get
	{
		fmt.Println("----Get----")
		r, err := cli.Get(context.Background(), "key").Result()
		fmt.Println("result", r)
		fmt.Println("err", err)
	}

	// Ping
	{
		fmt.Println("----Ping----")
		r, err := cli.Ping(context.Background()).Result()
		fmt.Println(r, err)
	}

	// HSet
	{
		fmt.Println("----HSet----")
		r, err := cli.HSet(context.Background(), "hset-key", "field", "1").Result()
		fmt.Println("result:", r)
		fmt.Println("err:", err)
	}

	// HGet
	{
		fmt.Println("----HGet----")
		r, err := cli.HGet(context.Background(), "hset-key", "field").Result()
		fmt.Println("result:", r)
		fmt.Println("err:", err)
	}
}
