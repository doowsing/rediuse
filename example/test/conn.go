package main

import (
	"github.com/gomodule/redigo/redis"
	"github.com/mna/redisc"
	"log"
	"time"
)

var cluster *redisc.Cluster

func init() {
	cluster = &redisc.Cluster{
		StartupNodes: []string{":7000", ":7001", ":7002"},
		DialOptions:  []redis.DialOption{redis.DialConnectTimeout(5 * time.Second)},
		CreatePool:   createPool,
	}
	cluster.Stats()
	if err := cluster.Refresh(); err != nil {
		log.Fatalf("Refresh failed: %v", err)
	}
}

func createPool(addr string, opts ...redis.DialOption) (*redis.Pool, error) {
	return &redis.Pool{
		MaxIdle:     500,
		MaxActive:   1000,
		IdleTimeout: time.Minute,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", addr, opts...)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}, nil
}

func GetRedisCluster() *redisc.Cluster {
	return cluster

}

func GetRedisConn() redis.Conn {
	return cluster.Get()
}
