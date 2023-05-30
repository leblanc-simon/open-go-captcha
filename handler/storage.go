package handler

import (
	"context"
	"encoding/json"
	"sort"
	"time"

	"github.com/redis/go-redis/v9"
	"leblanc.io/open-go-captcha/connection"
	"leblanc.io/open-go-captcha/log"
)

var ctx = context.Background()

func getKey(session string) string {
	return connection.GetRedisKeyPrefix() + session
}

func getValidKey(session string) string {
	return connection.GetRedisKeyPrefix() + "_valid_" + session
}

func SetAnswer(session string, results []string) (bool) {
	rdb := connection.GetRedisInstance()

	sort.Strings(results)
	resultsToStore, _ := json.Marshal(results)

	err := rdb.Set(
		ctx,
		getKey(session),
		resultsToStore,
		time.Duration(connection.GetRedisExpire() * int(time.Second)),
	).Err()

	if err != nil {
		log.Error(err.Error())

		return false
	}

	return true
}

func StoreValidResult(session string) error {
	rdb := connection.GetRedisInstance()

	return rdb.Set(
		ctx,
		getValidKey(session),
		"1",
		time.Duration(connection.GetRedisLongExpire() * int(time.Second)),
	).Err()
}

func CheckAnswer(session string, answers []string) (bool) {
	rdb := connection.GetRedisInstance()

	sort.Strings(answers)
	answersToCompare, _ := json.Marshal(answers)

	val, err := rdb.Get(ctx, getKey(session)).Result()
	if err == redis.Nil {
		log.Error("Fail to found " + session)
		return false
	}
	
	if err != nil {
		log.Error(err.Error())
		return false
	}

	// Delete key to avoid brute force
	err = rdb.Del(ctx, getKey(session)).Err()
	if err != nil {
		log.Error(err.Error())
		return false
	}

	if string(answersToCompare) != val {
		log.Error("answer is not valid")
		return false
	}

	return true
}

func CheckValidResult(session string) bool {
	rdb := connection.GetRedisInstance()

	val, err := rdb.Get(ctx, getValidKey(session)).Result()
	if err == redis.Nil {
		log.Error("Fail to found " + session)
		return false
	}
	
	if err != nil {
		log.Error(err.Error())
		return false
	}

	err = rdb.Del(ctx, getValidKey(session)).Err()
	if err != nil {
		log.Error(err.Error())
		return false
	}

	if val != "1" {
		log.Error("answer is not valid")
		return false
	}

	return true
}