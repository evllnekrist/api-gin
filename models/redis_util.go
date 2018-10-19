// command sort by data-type, ditambah berdasarkan --> https://redis.io/topics/data-types
package models

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

type RedisUtil struct{}

func (model_redis RedisUtil) RedisPing() error {

	conn := Pool.Get()
	defer conn.Close()

	_, err := redis.String(conn.Do("PING"))
	if err != nil {
		return fmt.Errorf("cannot 'PING' db: %v", err)
	} else {
		fmt.Sprintf("Pinging %s.. successfully", myRedisHost)
	}
	return nil
}

// STRING start
func (model_redis RedisUtil) RedisGet(key string) (string, bool) {

	conn := Pool.Get()
	defer conn.Close()

	var data string
	data, err := redis.String(conn.Do("GET", key))

	if err != nil {
		return data, false
	}
	return data, true
}

func (model_redis RedisUtil) RedisSet(key string, value []byte) error {

	conn := Pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, value)
	if err != nil {
		v := string(value)
		if len(v) > 15 {
			v = v[0:12] + "..."
		}
		return fmt.Errorf("error setting key %s to %s: %v", key, v, err)
	}
	return err
}

func (model_redis RedisUtil) RedisExists(key string) (bool, error) {

	conn := Pool.Get()
	defer conn.Close()

	ok, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return ok, fmt.Errorf("error checking if key %s exists: %v", key, err)
	}
	return ok, err
}

func (model_redis RedisUtil) RedisDelete(key string) error {

	conn := Pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", key)
	return err
}

func (model_redis RedisUtil) RedisGetKeys(pattern string) ([]string, error) {

	conn := Pool.Get()
	defer conn.Close()

	iter := 0
	keys := []string{}
	for {
		arr, err := redis.Values(conn.Do("SCAN", iter, "MATCH", pattern))
		if err != nil {
			return keys, fmt.Errorf("error retrieving '%s' keys", pattern)
		}

		iter, _ = redis.Int(arr[0], nil)
		k, _ := redis.Strings(arr[1], nil)
		keys = append(keys, k...)

		if iter == 0 {
			break
		}
	}

	return keys, nil
}

func (model_redis RedisUtil) RedisIncr(counterKey string) (int, error) {

	conn := Pool.Get()
	defer conn.Close()

	return redis.Int(conn.Do("INCR", counterKey))
}

//LIST start
func (model_redis RedisUtil) RedisLrange(key string, limit string, start string) ([]string, bool) { //get ouput

	conn := Pool.Get()
	defer conn.Close()

	data, err := redis.Strings(conn.Do("LRANGE", key, start, limit))
	if err != nil {
		return data, false
	}
	return data, true
}

func (model_redis RedisUtil) RedisRpush(key string, value []byte) error { //do input

	conn := Pool.Get()
	defer conn.Close()

	_, err := conn.Do("RPUSH", key, value)
	if err != nil {
		v := string(value)
		if len(v) > 15 {
			v = v[0:12] + "..."
		}
		return fmt.Errorf("error setting key %s to %s: %v", key, v, err)
	}
	return err
}

func (model_redis RedisUtil) RedisLlen(counterKey string) (int, error) {

	conn := Pool.Get()
	defer conn.Close()

	return redis.Int(conn.Do("LLEN", counterKey))
}

//SET start
func (model_redis RedisUtil) RedisSmembers(key string, limit string, start string) ([]string, bool) { //get ouput

	conn := Pool.Get()
	defer conn.Close()

	data, err := redis.Strings(conn.Do("SMEMBERS", key, start, limit))
	if err != nil {
		return data, false
	}
	return data, true
}

func (model_redis RedisUtil) RedisSadd(key string, value []byte) error { //do input

	conn := Pool.Get()
	defer conn.Close()

	_, err := conn.Do("SADD", key, value)
	if err != nil {
		v := string(value)
		if len(v) > 15 {
			v = v[0:12] + "..."
		}
		return fmt.Errorf("error setting key %s to %s: %v", key, v, err)
	}
	return err
}

//ZSET start UNDEFINED

//HASH start UNDEFINED

//oth.
func (model_redis RedisUtil) RedisKeys(counterKey string) ([]string, error) {

	conn := Pool.Get()
	defer conn.Close()

	counterKey = counterKey + "*"
	return redis.Strings(conn.Do("KEYS", counterKey))
}

func (model_redis RedisUtil) RedisType(counterKey string) (string, error) {

	conn := Pool.Get()
	defer conn.Close()

	return redis.String(conn.Do("TYPE", counterKey))
}
