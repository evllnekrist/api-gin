package models

import (
	"api-gin/helpers"
	"github.com/gomodule/redigo/redis"

	"os"
	"os/signal"
	"syscall"
	"time"
)

type RedisPool struct{}

var helper_eve = new(helpers.EveHelper) //nama-package.nama-type-data-file-package

var (
	Pool        *redis.Pool
	myRedisHost = "172.17.6.45"
	mydb_redis  = "1"
	// myRedisHost = "192.168.4.191"
	// mydb_redis   = "0"
)

func (db_redis RedisPool) Init() {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = myRedisHost + ":6379"
	}
	Pool = db_redis.newPool(redisHost)
	err := db_redis.selectDb(mydb_redis)
	helper_eve.Panics(err)
	db_redis.cleanupHook()
}

func (db_redis RedisPool) newPool(server string) *redis.Pool {

	return &redis.Pool{

		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,

		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			return c, err
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func (db_redis RedisPool) selectDb(key string) error {
	conn := Pool.Get()
	defer conn.Close()

	_, err := conn.Do("SELECT", key)
	return err
}

func (db_redis RedisPool) cleanupHook() {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGKILL)
	go func() {
		<-c
		Pool.Close()
		os.Exit(0)
	}()
}
