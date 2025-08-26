package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
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

// HashCRUD
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

	// ======== HGET: 获取所有字段的值 ========
	profile, _ := rdb.HGetAll(ctx, hashKey).Result()
	fmt.Println("HGETALL获取key的所有信息:", profile)
	fmt.Printf("%s的类型是:%v\n", "profile", reflect.TypeOf(profile))

	// ======== HINCRBY: 对字段值进行原子增减 ========
	// 应用场景:统计用户积分、登录次数(可避免先读取再写入的复杂操作)
	// 1为原子增,-1为原子减
	newAge, _ := rdb.HIncrBy(ctx, hashKey, "age", 1).Result()
	fmt.Printf("key:%s,HINCRBY修改后的%s字段值为:%v\n", hashKey, "age", newAge)
	fmt.Printf("修改后新值的类型是:%v\n", reflect.TypeOf(newAge))

	// ======== HEXISTS: 检查字段是否存在 ========
	exists, _ := rdb.HExists(ctx, hashKey, "city").Result()
	fmt.Printf("字段'%s' 是否存在:%v\n", "city", exists)

	// ======== HDEL: 删除一个或多个字段 ========
	rdb.HDel(ctx, hashKey, "email")
	fmt.Println("HDel 删除 email 字段后")
	profileAfterDel, _ := rdb.HGetAll(ctx, hashKey).Result()
	fmt.Println("当前用户信息:", profileAfterDel)

}

// SetCRUD
func SetOps() {
	setKey1 := "article:1:tags" // 文章1标签
	setKey2 := "article:2:tags" // 文章2标签
	clearKeys(setKey1)
	clearKeys(setKey2)

	// ======== SADD: 添加一个或多个成员 ========
	// 应用场景：为文章、商品等打标签
	rdb.SAdd(ctx, setKey1, "go", "redis", "web")
	rdb.SAdd(ctx, setKey2, "go", "docekr", "performance")
	fmt.Println("为两个文章添加标签")

	// ======== SMEMBERS: 获取集合中的所有成员 =======
	tags1, _ := rdb.SMembers(ctx, setKey1).Result()
	fmt.Println("文章1的所有标签:", tags1)

	// ======== SISMEMBER: 判断成员是否存在于集合中 ========
	isMember, _ := rdb.SIsMember(ctx, setKey1, "redis").Result()
	fmt.Printf("'%s'是不是文章1的标签:%v\n", "redis", isMember)

	// ======== SCARD: 获取集合的成员数量 =======
	count, _ := rdb.SCard(ctx, setKey2).Result()
	fmt.Printf("文章'%s'标签的数量是:%v\n", setKey2, count)

	// ======== SINTER: 获取两个或多个集合的交集 ========
	// 应用场景：发现共同兴趣、共同好友、相关文章等
	commonTags, _ := rdb.SInter(ctx, setKey1, setKey2).Result()
	fmt.Println("两篇文章的共同标签(交集):", commonTags)

	// ======== SUNION: 获取两个或多个集合的并集 ========
	allTags, _ := rdb.SUnion(ctx, setKey1, setKey2).Result()
	fmt.Println("两篇文章的所有标签(并集):", allTags)

	// ======== SREM: 从集合中移除一个或多个成员 ========
	rdb.SRem(ctx, setKey1, "web")
	tags1AfterRem, _ := rdb.SMembers(ctx, setKey1).Result()
	fmt.Println("从文章1移除 'web' 标签后:", tags1AfterRem)
}

// SortedCRUD
func SortedSetOps() {
	// Zset存储的也是字符串成员类型；每个成员关联一个分数(score)；大多数操作时间复杂度是O(logN)
	// 排行榜系统、新闻Feed
	// 与Hash配合使用，先查id排名，再查id信息
	zsetKey := "game:leaderboard"
	clearKeys(zsetKey)

	// ======== ZADD: 添加一个或多个成员，每个成员都有一个分数(score) ========
	// 应用场景：排行榜、带权重的任务队列等
	rdb.ZAdd(ctx, zsetKey, redis.Z{Score: 1500, Member: "PlayerOne"})
	rdb.ZAdd(ctx, zsetKey, redis.Z{Score: 1400, Member: "PlayerTwo"})
	rdb.ZAdd(ctx, zsetKey, redis.Z{Score: 1300, Member: "PlayerThree"})
	fmt.Println("添加三条玩家数据")

	// ======== ZSCORE: 获取指定成员的分数 ========
	score, _ := rdb.ZScore(ctx, zsetKey, "PlayerOne").Result()
	fmt.Println("PlayerOne 的当前分数:", score)

	// ======== ZRANK/ZREVRANK: 获取成员的排名 (升序/降序) ========
	// 排名从0开始
	rank, _ := rdb.ZRevRank(ctx, zsetKey, "PlayerOne").Result() // ZREVRANK 用于从高到低的排名
	fmt.Println("PlayerOne 的当前排名(从高到低):", rank)

	// ======== ZRANGE/ZREVRANGE: 按排名范围获取成员 ========
	// 应用场景：获取排行榜Top N
	// ZRevRangeWithScores: 获取指定范围的成员，并带上他们的分数 (从高到低)
	top3, _ := rdb.ZRevRangeWithScores(ctx, zsetKey, 0, 2).Result() // 获取排名前3的玩家 (0, 1, 2)
	fmt.Println("排行榜 Top 3:")
	for _, player := range top3 {
		fmt.Printf("  - 玩家: %s, 分数: %s\n", player.Member, strconv.FormatFloat(player.Score, 'f', -1, 64))
	}

	// ======== ZCOUNT: 获取指定分数区间的成员数量 ========
	count, _ := rdb.ZCount(ctx, zsetKey, "2000", "3500").Result() // 分数在 [2000, 3500] 之间的玩家数量
	fmt.Println("分数在2000到3500之间的玩家数量:", count)
}

func main() {
	initRedis()
	StringOps()
	ListOps()
	HashOps()
	SetOps()
	SortedSetOps()
}
