package connection

import (
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
	"leblanc.io/open-go-captcha/config"
)

var lockRedis = &sync.Mutex{}

type redisStruct struct {
	Options *redis.Options
	Conn *redis.Client
	Expire int
	LongExpire int
	KeyPrefix string
}

var redisInstance *redisStruct

func Initialize(c *config.Config) {
	lockRedis.Lock()
	defer lockRedis.Unlock()

	redisInstance = &redisStruct{}
	redisInstance.Options = &redis.Options{
		Addr: fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port),
		Password: c.Redis.Password,
		DB: c.Redis.Db,
	}

	redisInstance.Conn = redis.NewClient(redisInstance.Options)
	redisInstance.Expire = c.Redis.Expire
	redisInstance.LongExpire = c.Redis.LongExpire
	redisInstance.KeyPrefix = c.Redis.KeyPrefix
}

func GetRedisInstance() *redis.Client {
	return redisInstance.Conn
}

func GetRedisExpire() int {
	return redisInstance.Expire
}

func GetRedisLongExpire() int {
	return redisInstance.LongExpire
}

func GetRedisKeyPrefix() string {
	return redisInstance.KeyPrefix
}