package rediuse

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/unknwon/com"
	"reflect"
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
	result.args = args
	result.SetResult(conn.Do(commandName, args...))
	return result
}

func (handler *RdbHandler) Get(key string) *RdbResult {
	return handler.DoCommand("GET", key)
}

func (handler *RdbHandler) Set(key string, value interface{}) *RdbResult {
	value, err := marshal(value)
	if err != nil {
		return &RdbResult{err: err}
	}
	return handler.DoCommand("SET", key, value)
}

func (handler *RdbHandler) SetEx(key string, value interface{}, second int) *RdbResult {
	value, err := marshal(value)
	if err != nil {
		return &RdbResult{err: err}
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
		return &RdbResult{err: err}
	}
	return handler.DoCommand("HSET", key, getString(field), value)
}

func (handler *RdbHandler) Hdel(key string, field interface{}) *RdbResult {
	return handler.DoCommand("HDEL", key, getString(field))
}

func (handler *RdbHandler) Hexist(key string, field interface{}) *RdbResult {
	return handler.DoCommand("HEXIST", key, getString(field))
}

func (handler *RdbHandler) HgetAll(key string) RdbHashResultInterface {
	return handler.DoCommand("HGETALL", key)
}

func (handler *RdbHandler) Hmget(key string, fields ...interface{}) RdbHashResultInterface {
	return handler.DoCommand("HMGET", redis.Args{}.Add(key).AddFlat(fields)...)
}

func (handler *RdbHandler) Hmset(key string, filed2data interface{}) *RdbResult {

	// 第一种方法，但是这种方法没有过滤 传参的类型
	//args := redis.Args{}.AddFlat(key).AddFlat(filed2data)

	// 第二种方法，复制redis.Args{}..AddFlat()方法，并过滤参数类型
	rv := reflect.ValueOf(filed2data)
	var args = []interface{}{key}
	switch rv.Kind() {
	case reflect.Slice:
		rvLen := rv.Len()
		if rvLen%2 == 1 {
			rvLen--
		}
		var isKey = true
		for i := 0; i < rvLen; i++ {
			if isKey {
				args = append(args, rv.Index(i).Interface())
			} else {
				vJson, err := marshal(rv.Index(i).Interface())
				if err != nil {
					return NewErrRdbResult("批量序列化失败")
				} else {
					args = append(args, vJson)
				}
			}
			isKey = !isKey

		}
	case reflect.Map:
		for _, k := range rv.MapKeys() {
			vJson, err := marshal(rv.MapIndex(k).Interface())
			if err != nil {
				return NewErrRdbResult("批量序列化失败")
			} else {
				args = append(args, getString(k.Interface()), vJson)
			}
		}
	default:
		return NewErrRdbResult("非 map or slice，无法写入redis hash队列")
	}

	// 第三种方法，限定了redis只能传 map[string]interface
	//var args = []interface{}{key}
	//for i, v := range filed2data {
	//	vJson, err := marshal(v)
	//	if err != nil {
	//		return NewErrRdbResult("批量序列化失败")
	//	} else {
	//		args = append(args, i, vJson)
	//	}
	//}
	//fmt.Printf("hmset args:%v\n", args)
	return handler.DoCommand("HMSET", args...)
}

func (handler *RdbHandler) RPush(key string, field interface{}) *RdbResult {
	return handler.DoCommand("RPUSH", key, getString(field))
}

func (handler *RdbHandler) LPush(key string, field interface{}) *RdbResult {
	return handler.DoCommand("LPUSH", key, getString(field))
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

func (handler *RdbHandler) HScan(key string, cursor int, match string, count int) (int, map[string]string, error) {
	var args []interface{}
	args = append(args, key, cursor)
	if match != "" {
		args = append(args, "match", match)
	}
	if count > 0 {
		args = append(args, "count", count)
	}

	result, err := handler.DoCommand("HSCAN", args...).Interface()
	mapResult := make(map[string]string)
	if err == nil {
		_datas := result.([]interface{})
		nextCursor := string(_datas[0].([]byte))
		datas := _datas[1].([]interface{})
		cursor = com.StrTo(nextCursor).MustInt()
		//fmt.Println(nextCursor)
		for i := 0; i < len(datas)/2; i++ {
			mapResult[string(datas[i*2].([]byte))] = mapResult[string(datas[i*2+1].([]byte))]
		}
	}
	return cursor, mapResult, err
}

func (handler *RdbHandler) Scan(cursor int, match string, count int) (int, []string, error) {
	var args []interface{}
	args = append(args, cursor)
	if match != "" {
		args = append(args, "match", match)
	}
	if count > 0 {
		args = append(args, "count", count)
	}
	result, err := handler.DoCommand("SCAN", args...).Interface()
	mapResult := []string{}
	if err == nil {
		_datas := result.([]interface{})
		nextCursor := string(_datas[0].([]byte))
		datas := _datas[1].([]interface{})
		//fmt.Println(nextCursor)
		cursor, err = com.StrTo(nextCursor).Int()
		if err != nil {
			fmt.Printf("返回游标非整数，%s\n", nextCursor)
		}
		for i := 0; i < len(datas); i++ {
			mapResult = append(mapResult, string(datas[i].([]byte)))
		}
	}
	return cursor, mapResult, err
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
