# redisx

封装了 go-redis，支持单机和集群

配置

```go
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
```

#### 安装

`go get -u -v github.com/itmisx/redisx`

#### 使用

参考 redix_test.go，详细请查看 github: https://github.com/go-redis/redis
