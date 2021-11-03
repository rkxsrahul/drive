package redisdb

import (
	"fmt"
	"log"
	"strconv"

	"git.xenonstack.com/util/drive-portal/config"
	"github.com/go-redis/redis"
)

var Client *redis.Client

func Initialise() error {
	db, _ := strconv.Atoi(config.RedisDB)
	// redis db client creations
	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort),
		Password: fmt.Sprintf("%s", config.RedisPass),
		DB:       db,
	})

	// check connection with server
	pong, err := Client.Ping().Result()
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(pong)
	return nil
}

// fucntion to checktoken exists
func CheckToken(token string) error {
	_, err := Client.Get(token).Result()
	if err != nil {
		// when token not exist
		log.Println(err)
		return err
	}
	// log.Println(val)
	return nil
}

// function for saving token in redis
func SaveToken(key, value string) {
	err := Client.Set(key, value, config.JWTExpireTime).Err()
	if err != nil {
		panic(err)
	}
}

// function for deleting token from redis
func DeleteToken(key string) error {
	val, err := Client.Del(key).Result()
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(val)
	return nil
}
