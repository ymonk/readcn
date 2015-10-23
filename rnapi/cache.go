package main

import (
    "gopkg.in/redis.v3"
)

const (
    KEY = "readcn_web_cache"
)

type Cache struct {
    redis *redis.Client
}

func NewCache() *Cache {
    rc := redis.NewClient(
        &redis.Options{
            Addr: "localhost:6379",
            Password: "",
            DB: 2,
        },
    )
    return &Cache{redis: rc}
}

func (c *Cache) hit(key string) bool {
    result, err := c.redis.HExists(KEY, key).Result()
    if err == nil && result == true {
        return true
    }
    return false
}

func (c *Cache) retrieve(key string) ([]byte, error) {
    return c.redis.HGet(KEY, key).Bytes()
}

func (c *Cache) store(key string, bytes []byte) {
    c.redis.HSet(KEY, key, string(bytes))
}

