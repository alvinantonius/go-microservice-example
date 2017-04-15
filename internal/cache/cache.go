package cache

import (
	"fmt"

	"gopkg.in/redis.v5"
)

var connections map[string]*redis.Client

// ConnectRedis for initiate connection to redis server, and store it as private global variable
func ConnectRedis(config map[string]string) {

	connections = make(map[string]*redis.Client)

	for name, addr := range config {
		conn := redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: "",
		})

		connections[name] = conn
	}
}

// Conn is for get redis connection
func Conn(name string) (*redis.Client, error) {
	if rds, ok := connections[name]; ok {
		return rds, nil
	}

	return nil, fmt.Errorf("redis conn not found")
}
