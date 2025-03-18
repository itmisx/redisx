package redisx

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestRedisx(*testing.T) {
	// 配置
	config := Config{
		Host:     "127.0.0.1",
		Port:     "16379",
		Protocol: "tcp",
		Prefix:   "prod:",
	}
	// 新建实例
	rdb := New(config)
	// 命令测试

	// // SetEx
	// {
	// 	fmt.Println("----Set----")
	// 	r, err := cli.SetEX(context.Background(), "key", "123", time.Second*10).Result()
	// 	fmt.Println("result", r)
	// 	fmt.Println("err", err)
	// }

	// Set
	{
		fmt.Println("----Set1----")
		r, err := rdb.Set(context.Background(), "key1", 1, time.Second*60).Result()
		fmt.Println("result", r)
		fmt.Println("err", err)

		fmt.Println("----Set2----")
		r, err = rdb.Set(context.Background(), "key2", 1, time.Second*60).Result()
		fmt.Println("result", r)
		fmt.Println("err", err)
	}
	// exits
	{
		fmt.Println("----Exits----")
		count, err := rdb.Exists(context.Background(), "key1", "key2").Result()
		fmt.Println("count", count)
		fmt.Println("err", err)
	}
	// pipeline
	{
		fmt.Println("----Pipeline----")
		pipe := rdb.Pipeline()
		pipe.Exists(context.Background(), "key1")
		pipe.Exists(context.Background(), "key2")
		cmder, err := pipe.Exec(context.Background())
		if err != nil {
			fmt.Println("Pipeline Exec Error:", err)
		}
		for _, cmd := range cmder {
			fmt.Println(cmd.String())
		}
	}

	// Del
	// {
	// 	fmt.Println("----Del---")
	// 	r, err := cli.Del(context.Background(), "key1", "key2").Result()
	// 	fmt.Println("result", r)
	// 	fmt.Println("err", err)
	// }
	// Get
	// {
	// 	fmt.Println("----Get----")
	// 	r, err := cli.Get(context.Background(), "key").Result()
	// 	fmt.Println("result", r)
	// 	fmt.Println("err", err)
	// }

	// // Ping
	// {
	// 	fmt.Println("----Ping----")
	// 	r, err := cli.Ping(context.Background()).Result()
	// 	fmt.Println(r, err)
	// }

	// // HSet
	// {
	// 	fmt.Println("----HSet----")
	// 	r, err := cli.HSet(context.Background(), "hset-key", "field", "1").Result()
	// 	fmt.Println("result:", r)
	// 	fmt.Println("err:", err)
	// }

	// // HGet
	// {
	// 	fmt.Println("----HGet----")
	// 	r, err := cli.HGet(context.Background(), "hset-key", "field").Result()
	// 	fmt.Println("result:", r)
	// 	fmt.Println("err:", err)
	// }
}
