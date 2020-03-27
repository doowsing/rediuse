package rediuse

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"strconv"
)

type RdbHandler struct {
	getConn func() redis.Conn
}

func NewRdbHandler(getConn func() redis.Conn) *RdbHandler {
	return &RdbHandler{getConn: getConn}
}

func (handler *RdbHandler) DoCommand(commandName string, args ...interface{}) *RdbResult {
	conn := handler.getConn()
	defer conn.Close()
	result := &RdbResult{}
	result.SetResult(conn.Do(commandName, args...))
	return result
}

func (handler *RdbHandler) Get(key string) *RdbResult {
	return handler.DoCommand("GET", key)
}

func (handler *RdbHandler) Set(key string, value interface{}) *RdbResult {
	value, err := marshal(value)
	if err != nil {
		return NewErrRdbResult(err)
	}
	return handler.DoCommand("SET", key, value)
}

func (handler *RdbHandler) SetEx(key string, value interface{}, second int) *RdbResult {
	value, err := marshal(value)
	if err != nil {
		return NewErrRdbResult(err)
	}
	return handler.DoCommand("SET", key, value, "EX", second)
}

func (handler *RdbHandler) Expire(key string, second int) *RdbResult {
	return handler.DoCommand("EXPIRE", key, second)
}

func (handler *RdbHandler) Exists(key string) *RdbResult {
	return handler.DoCommand("Exists", key)
}

func (handler *RdbHandler) Delete(key string) *RdbResult {
	return handler.DoCommand("DEL", key)
}

func (handler *RdbHandler) Hget(key string, field interface{}) *RdbResult {
	return handler.DoCommand("HGET", key, getString(field))
}

func (handler *RdbHandler) Hset(key string, field interface{}, value interface{}) *RdbResult {
	value, err := marshal(value)
	if err != nil {
		return NewErrRdbResult(err)
	}
	return handler.DoCommand("HSET", key, getString(field), value)
}

func (handler *RdbHandler) Hdel(key string, field interface{}) *RdbResult {
	return handler.DoCommand("HDEL", key, getString(field))
}

func (handler *RdbHandler) Hexist(key string, field interface{}) *RdbResult {
	return handler.DoCommand("HEXIST", key, getString(field))
}

func (handler *RdbHandler) Hmget(key string) *RdbResult {
	return handler.DoCommand("HMGET", key)
}

func (handler *RdbHandler) Hmset(key string, filed2data map[string]interface{}) *RdbResult {
	var args = []interface{}{key}
	for i, v := range filed2data {
		vJson, err := marshal(v)
		if err != nil {
			return NewErrRdbResult(err)
		} else {
			args = append(args, i, vJson)
		}
	}
	return handler.DoCommand("HMGET", args...)
}

func (handler *RdbHandler) RPush(key string, field interface{}) *RdbResult {
	return handler.DoCommand("RPUSH", key, getString(field))
}

func (handler *RdbHandler) LLen(key string) *RdbResult {
	return handler.DoCommand("LLEN", key)
}

func (handler *RdbHandler) LRange(key string, start, stop int) *RdbResult {
	return handler.DoCommand("LRANGE", key, start, stop)
}

func (handler *RdbHandler) LTrim(key string, start, stop int) *RdbResult {
	return handler.DoCommand("LTRIM", key, start, stop)
}

func (handler *RdbHandler) SCARD(key string) *RdbResult {
	return handler.DoCommand("SCARD", key)
}

func (handler *RdbHandler) SADD(key string, field interface{}) *RdbResult {
	return handler.DoCommand("SADD", key, getString(field))
}

func getString(field interface{}) string {
	trueField := ""
	switch field.(type) {
	case int:
		trueField = strconv.Itoa(field.(int))
		break
	case string:
		trueField = field.(string)
		break
	default:
		trueField = fmt.Sprintf("%s", field)
	}
	return trueField
}

func marshal(v interface{}) (interface{}, error) {
	switch t := v.(type) {
	case string, []byte, int, int64, float64, bool, nil, redis.Argument:
		return v, nil
	default:
		return json.Marshal(t)
	}
}
