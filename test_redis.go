package main

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

func main() {
	// Create Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	ctx := context.Background()

	// Test PING
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Redis connection failed:", err)
	}
	fmt.Println("Redis PING:", pong)

	// Test SET
	err = client.Set(ctx, "test_key", "test_value", 0).Err()
	if err != nil {
		log.Fatal("Redis SET failed:", err)
	}
	fmt.Println("Redis SET: OK")

	// Test GET
	val, err := client.Get(ctx, "test_key").Result()
	if err != nil {
		log.Fatal("Redis GET failed:", err)
	}
	fmt.Println("Redis GET:", val)

	fmt.Println("\n✅ Redis connection is working!")
}
