package main

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

var rdb *redis.Client = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKMNOPQRSTUVWXYZ")

// RandomString 随机生成一个长度为 size 的字符串
func RandomString(size int) string {
	b := make([]rune, size)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// GetUsedMemory 获取Redis已经使用的内存大小
func GetUsedMemory() (int, error) {
	status, err := rdb.Do("INFO", "MEMORY").Result()
	if err != nil {
		return -1, errors.Wrap(err, "获取Redis内存失败")
	}
	re := regexp.MustCompile(`used_memory:(\d+)`)
	matched := re.FindStringSubmatch(status.(string))
	return strconv.Atoi(matched[1])
}

const count = 10000

// Write 写入10000个数据到数据库
func Write(size int) {
	for i := 0; i < count; i++ {
		key := fmt.Sprintf("%d:%d", size, i)
		value := RandomString(size)
		rdb.Set(key, value, time.Minute*20)
	}
}

func main() {
	sizes := []int{5, 10, 20, 50, 100, 1000, 10000}
	fmt.Printf("写入数量为: %v\n", count)
	for _, size := range sizes {
		rdb.FlushAll()

		beforeMemory, err := GetUsedMemory()
		if err != nil {
			fmt.Printf("写入前获取内存占用失败:%v\n", err)
			panic(err)
		}
		fmt.Printf("写入前内存大小为: %v\n", beforeMemory)

		Write(size)

		afterMomory, err := GetUsedMemory()
		if err != nil {
			fmt.Printf("写入后获取内存占用失败:%v\n", err)
			panic(err)
		}
		fmt.Printf("写入后内存大小为: %v\n", afterMomory)

		fmt.Printf("大小为%d的数据的平均内存为:%d\n", size, (afterMomory-beforeMemory)/count)
	}
}
