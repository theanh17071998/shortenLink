package redis

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

//Pool Kết nối đến redis
var Pool *redis.Pool

//Hàm tạo connection pool đến redis
func newPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     80,
		MaxActive:   12000, // số lượng kết nối max
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "redis-master.shorten-url.svc.cluster.local:6379")
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}

// Ping Kiểm tra kết nối
func Ping() error {
	conn := Pool.Get()
	defer conn.Close()
	s, err := redis.String(conn.Do("PING"))
	if err != nil {
		return err
	}
	fmt.Printf("PING Response = %s\n", s)
	return nil
}

//Hàm khởi tạo
func init() {
	Pool = newPool()
	err := Ping()
	if err != nil {
		panic(err)
	}
	return
}

// Set Cache lại shortcode và longurl trả về lỗi nếu có
func Set(shortcode, longurl string) error {
	conn := Pool.Get()
	defer conn.Close()
	_, err := conn.Do("SET", shortcode, longurl)
	return err
}

//Get Lấy longurl từ cache trả về lỗi nếu có
func Get(shortcode string) (string, error) {
	conn := Pool.Get()
	defer conn.Close()
	longurl, err := redis.String(conn.Do("GET", shortcode))
	if err == redis.ErrNil {
		return "", nil
	}
	return longurl, err
}
