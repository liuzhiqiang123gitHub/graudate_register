package redisUtil

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	config "graduate_registrator/utils/conf"
)

var (
	RedisConn redis.Conn
)

func InitRedis(redisConf config.Configure) (err error) {
	//fmt.Printf(" %+v", redisConf)
	RedisConn, err = redis.Dial("tcp", redisConf.RedisSetting["tmp"].RedisConn, redis.DialPassword(redisConf.RedisSetting["tmp"].RedisPasswd))
	if err != nil {
		fmt.Println("Connect to redis error", err)
	}
	return err
}

//插入数据
func Set(key, value interface{}, exp uint) error {
	_, err := RedisConn.Do("SET", key, value, "EX", exp)
	if err != nil {
		fmt.Println("redis set failed:", err)
		return err
	}
	return nil
}

//Get
func Get(key interface{}) (resStr interface{}, err error) {
	res, err := redis.String(RedisConn.Do("GET", key))
	return res, err
}

//Del
func Delete(key string) error {
	_, err := RedisConn.Do("DEL", key)
	return err
}
