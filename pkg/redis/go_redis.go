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

// ClearKeys
func clearKeys(prefix string) {
	iter := rdb.Scan(ctx, 0, prefix+"*", 0).Iterator()
	for iter.Next(ctx) {
		err := rdb.Del(ctx, iter.Val()).Err()
		if err != nil {
			panic(err)
		}
	}
	if err := iter.Err(); err != nil {
		panic(err)
	}
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

// ListCRUD
func ListOps() {
	listKey := "tasks:queue"
	clearKeys(listKey)

	// ======== LPUSH/RPUSH: 从左/右堆入元素 ========
	rdb.LPush(ctx, listKey, "task_3")
	rdb.LPush(ctx, listKey, "task_2")
	rdb.RPush(ctx, listKey, "task_4")
	rdb.LPush(ctx, listKey, "task_1")
	fmt.Println("推入4个任务到队列")

	// ======== 获取列表长度 ========
	length, _ := rdb.LLen(ctx, listKey).Result()
	fmt.Println("当前任务队列长度:", length)

	// ======== LRANGE:获取制定范围的元素 ========
	tasks, _ := rdb.LRange(ctx, listKey, 0, -1).Result()
	fmt.Println("当前队列所有任务:", tasks)

	// ======== LPOP/RPOP 从左/右弹出元素 ========
	task, _ := rdb.RPop(ctx, listKey).Result()
	fmt.Println("处理了一个任务:", task)

	// ======== 阻塞式弹出BRPOP ========
	// 当列表为空时，BRPOP会阻塞等待，直到有新元素被推入或超时
	clearKeys(listKey)

	go func() {
		time.Sleep(3 * time.Second)
		fmt.Println("模拟[生产者]3s后推送数据")
		rdb.LPush(ctx, listKey, "urgent_task")
	}()

	// listKey对应的value为空，等待生产者推。设置超时时间为5s
	fmt.Println("[消费者]正在等待任务")
	result, err := rdb.BRPop(ctx, 5*time.Second, listKey).Result()
	if err != nil {
		fmt.Println("等待任务超时或出错:", err)
	} else {
		// result 是一个数组 [key, value]
		fmt.Printf("[消费者] 成功接收并处理任务: %s\n", result[1])
	}
}

// HaspCRUD
func HashOps() {
	hashKey := "user:2:profile"
	clearKeys(hashKey)

	// ======== HSET: 设置单个或多个字段 ========
	rdb.HSet(ctx, hashKey, "name", "RT", "age", "18", "email", "xxx@123") // 适合存储对象
	fmt.Println("使用HSet 设置用户信息")

	// ======== HGET: 获取单个字段的值 ========
	name, _ := rdb.HGet(ctx, hashKey, "name").Result()
	fmt.Printf("%skey的name字段的值是%s \n", hashKey, name)

	// ======== HGET: 获取多个字段的值 ========
	fields := []string{"name", "age"}
	fields_value, err := rdb.HMGet(ctx, hashKey, fields...).Result()
	if err != nil {
		log.Fatalf("HMget 出错: %v", err)
	}
	for i, field := range fields {
		fmt.Printf("%skey的%s字段的值是 %v \n", hashKey, field, fields_value[i])
	}
}

func main() {
	initRedis()
	// StringOps()
	ListOps()
}
