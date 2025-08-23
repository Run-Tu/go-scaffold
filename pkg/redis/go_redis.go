package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	rdb *redis.Client
	ctx = context.Background()
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func initRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Redis连接失败: %v", err)
	}
	fmt.Println("Redis连接成功")
}

// String CRUD
func StringOps() {
	// ======== String ========
	err := rdb.Set(ctx, "userinfo:string:1:name", "RT", 0).Err()
	if err != nil {
		log.Fatalf("Set Error: %v", err)
	}
	val, _ := rdb.Get(ctx, "userinfo:string:1:name").Result()
	fmt.Println("String", val)

	// ======== JSON Object ========
	user := User{ID: 1, Name: "RT", Email: "xxx@123.com"}
	jsonData, _ := json.Marshal(&user)
	err = rdb.Set(ctx, "user:1:info", jsonData, 0).Err()
	if err != nil {
		log.Fatalf("Set Json Object Error: %v", err)
	}
	val2, _ := rdb.Get(ctx, "user:1:info").Result()
	var user2 User
	unmashal_err := json.Unmarshal([]byte(val2), &user2)
	if unmashal_err != nil {
		log.Fatalf("Json Obejct Unmarshal Error: %v", err)
	}
	fmt.Println("JSON 对象:", user2.Email)

	// ======== With expiration ========
	err = rdb.Set(ctx, "session:1", "token-abc123", 10*time.Second).Err()
	if err != nil {
		log.Fatalf("Set expiration data Error: %v", err)
	}
	ttl, _ := rdb.TTL(ctx, "session:1").Result()
	fmt.Println("session:1(剩余TTL)", ttl)

	// ======== Byte Data ========
	imageData := []byte{0xFF, 0xD8, 0xFF} // 模拟JPEG文件头
	rdb.Set(ctx, "image:logo", imageData, 0)
	val3, _ := rdb.Get(ctx, "image:logo").Bytes()
	fmt.Println("二进制数据", val3)
}

func main() {
	initRedis()
	StringOps()
}
