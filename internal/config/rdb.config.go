package config

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RdbConfig struct {
	User string
	Pass string
	Host string
	Port string
}

func NewRdbConfig(user, pass, host, port string) *RdbConfig {
	return &RdbConfig{
		User: user,
		Pass: pass,
		Host: host,
		Port: port,
	}
}

func (rdb *RdbConfig) ConnectRdb() *redis.Client {
	connStr := fmt.Sprintf("%s:%s", rdb.Host, rdb.Port)

	return redis.NewClient(&redis.Options{
		Addr:     connStr,
		Username: rdb.User,
		Password: rdb.Pass,
	})
}
