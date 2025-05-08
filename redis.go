package redisx

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// Redis接口定义
// 因要兼容单例和集群，所以需定义接口，来兼容
// 接口定义了什么方法，外部才可以调用什么方法

// Redis 配置项
type Config struct {
	Addrs        []string `mapstructure:"addrs"`          // 连接地址。多个为集群
	MasterName   string   `mapstructure:"master_name"`    // 主节点名称，用于哨兵模式
	Password     string   `mapstructure:"password"`       // 密码
	DB           int      `mapstructure:"db"`             // 默认连接的数据库，仅支持单机模式
	MinIdleConns int      `mapstructure:"min_idle_conns"` // 最小空闲连接
	IdleTimeout  int      `mapstructure:"idle_timeout"`   // 空闲时间
	PoolSize     int      `mapstructure:"pool_size"`      // 连接池大小
	MaxConnAge   int      `mapstructure:"max_conn_age"`   // 连接最大可用时间
	Prefix       string   `mapstructure:"prefix"`         // 键前缀
}

// New 新建redis实例
func New(conf Config) redis.UniversalClient {
	// 默认闲置连接
	if conf.MinIdleConns == 0 {
		conf.MinIdleConns = 2
	}
	// 空闲超时时间，过期关闭空闲连接
	if conf.IdleTimeout == 0 || conf.IdleTimeout > 1800 {
		conf.IdleTimeout = 1800
	}
	// 默认连接池数量为2
	if conf.PoolSize == 0 {
		conf.PoolSize = 10
	}
	// 连接的生命周期为300秒
	if conf.MaxConnAge == 0 || conf.MaxConnAge > 3600 {
		conf.MaxConnAge = 3600
	}
	for {
		rdb := redis.NewUniversalClient(&redis.UniversalOptions{
			Addrs:           conf.Addrs,
			MasterName:      conf.MasterName,
			Password:        conf.Password,
			DB:              conf.DB,
			MinIdleConns:    conf.MinIdleConns,
			ConnMaxIdleTime: time.Second * time.Duration(conf.IdleTimeout),
			PoolSize:        conf.PoolSize,
			ConnMaxLifetime: time.Second * time.Duration(conf.MaxConnAge),
			DialTimeout:     time.Second * 5,
		})
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		res, err := rdb.Ping(ctx).Result()
		cancel()
		if strings.ToLower(res) != "pong" || err != nil {
			log.Println("redis connection failed,retry...")
			time.Sleep(time.Second * 5)
		} else {
			if conf.Prefix != "" {
				rdb.AddHook(&prefixHook{Prefix: conf.Prefix})
			}
			return rdb
		}
	}
}
