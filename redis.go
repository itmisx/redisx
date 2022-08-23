package redisx

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

// Redis接口定义
// 因要兼容单例和集群，所以需定义接口，来兼容
// 接口定义了什么方法，外部才可以调用什么方法
type Client interface {
	Subscribe(ctx context.Context, channels ...string) *redis.PubSub
	redis.Cmdable
	Close() error
}

// Redis 配置项
type Config struct {
	Cluster  bool   `mapstructure:"cluster" `
	Host     string `mapstructure:"host" `
	Port     string `mapstructure:"port" `
	Password string `mapstructure:"password"`
	Protocol string `mapstructure:"protocol"`
	Database int    `mapstructure:"database"`
	// 最小空闲连接
	MinIdleConns int `mapstructure:"min_idle_conns"`
	// 空闲时间
	IdleTimeout int `mapstructure:"idle_timeout"`
	// 连接池大小
	PoolSize int `mapstructure:"pool_size"`
	// 连接最大可用时间
	MaxConnAge int `mapstructure:"max_conn_age"`
}

// New 新建redis实例
func New(conf Config) Client {
	config := conf
	ctx := context.Background()
	hostMembers := strings.Split(config.Host, ",")

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

	if conf.Cluster || len(hostMembers) > 1 {
		for {
			// 集群
			rdb := redis.NewClusterClient(&redis.ClusterOptions{
				Addrs:        hostMembers,
				Password:     config.Password,
				MinIdleConns: config.MinIdleConns,
				IdleTimeout:  time.Second * time.Duration(config.IdleTimeout),
				PoolSize:     config.PoolSize,
				MaxConnAge:   time.Second * time.Duration(config.MaxConnAge),
			})
			res, err := rdb.Ping(ctx).Result()
			if strings.ToLower(res) != "pong" || err != nil {
				log.Println("redis connection failed,retry...")
				time.Sleep(time.Second * 5)
			} else {
				return rdb
			}
		}
	} else {
		for {
			rdb := redis.NewClient(&redis.Options{
				Addr:         config.Host + ":" + config.Port,
				Password:     config.Password,
				DB:           config.Database,
				MinIdleConns: config.MinIdleConns,
				IdleTimeout:  time.Second * time.Duration(config.IdleTimeout),
				PoolSize:     config.PoolSize,
				MaxConnAge:   time.Second * time.Duration(config.MaxConnAge),
			})
			res, err := rdb.Ping(ctx).Result()
			if strings.ToLower(res) != "pong" || err != nil {
				log.Println("redis connection failed,retry...")
				time.Sleep(time.Second * 5)
			} else {
				return rdb
			}
		}
	}
}
