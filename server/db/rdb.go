package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

type Client struct {
	*redis.Client
}

var (
	ErrNil = errors.New("no matching record found in redis database")
	Ctx    = context.TODO()
)

var once sync.Once
var client *Client

func RedisClient() {
	HOST := "localhost"
	PORT := "6379"
	USERNAME := "admin"
	PASS := "rahasia"

	once.Do(func() {
		conn := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", HOST, PORT),
			Username: USERNAME,
			Password: PASS,
		})

		_, err := conn.Ping(Ctx).Result()
		if err != nil {
			log.Fatalf("Could not connect to redis %v", err)
		}

		client = &Client{conn}
	})
	log.Println("Success connect to Redis")
}

func GetRedis() *Client {
	if client == nil {
		RedisClient()
	}
	return client
}

var allowValue int = 0
var allow *int = &allowValue

// var rcmdGet *redis.StringCmd

func CreateLimitter() (bool, string, error) {
	log.Println("allowed in createLimit", allowValue)
	limiterValue := client.TTL(context.Background(), "key-limiter")
	timeLeft := limiterValue.Val().Seconds()

	// handling if time exist and allow is 0
	if timeLeft >= 0 && allowValue == 0 {
		log.Println("Allowed request:", allowValue, "time:", timeLeft)
		return false, fmt.Sprintf("You can request after %v", timeLeft), nil
	}

	if allowValue != 0{
		// handling time expire
		if timeLeft <= 0{
			// create new limiter using TTL redis
			add, err := addLimiter()
			if err != nil{
				return add, "", err
			}
			return add, "", nil

		// handling if allow still exist
		} else {
			allowValue--
			if allowValue <= 0{
				// handling when allow value is 0 but time exist, after that condition in above
				if timeLeft >= 0{
					log.Println("Allowed request:", allowValue, "time:", timeLeft)
					return false, fmt.Sprintf("Your limit request is %v try after %v", allowValue, timeLeft), nil
				}
			}
			log.Println("Allowed:", allowValue, "time:", timeLeft)
		}
	} else {
		// create new limiter using TTL redis
		add, err := addLimiter()
		if err != nil{
			return add, "", err
		}
		return add, "", nil
	}

	// response format (status, handling message and error)
	return true, "", nil
}

func addLimiter() (bool, error){
	allowValue = 5
	ttl := time.Duration(1) * time.Minute
	key := "key-limiter"

	// store data using SET command
	rcmd := client.Set(context.Background(), key, allowValue, ttl)
	if err := rcmd.Err(); err != nil {
		fmt.Printf("unable to SET data. error: %v", err)
		return	false, err
	}

	// rcmdGet := client.Get(context.Background(), key)
	// if err := rcmdGet.Err(); err != nil {
	// 	fmt.Printf("unable to GET data. error: %v", err)
	// 	return false, err
	// }
	log.Println("Allow request: ", allowValue)
	return true, nil
}

// var rdb *redis.Client

// func InitRdb() {
// 	ctx := context.Background()
// 	rdb := redis.NewClient(&redis.Options{
// 		Addr: "localhost:6379",
// 		Username: "admin",
// 		Password: "rahasia",
// 	})
// 	_ = rdb.FlushDB(ctx).Err()
// 	// limiter := redis_rate.NewLimiter(rdb)
// 	// res, err := limiter.Allow(ctx, "limit_request", redis_rate.PerMinute(5))
// 	// if err != nil {
// 	// 	panic(err)
// 	// }
// 	// log.Println("allowed", res.Allowed, "remaining", res.Remaining, "time", res.Limit.Period)
// }

// func GetRdbConn() *redis.Client{
// 	return rdb
// }
// type Limitter struct {
// 	*redis_rate.Result
// }

// type AllowLimiter struct {
// 	*int
// }

// var res *Limitter

// var allow *AllowLimiter