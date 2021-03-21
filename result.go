package rediuse

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"reflect"
	"strconv"
)

type RdbResultInterface interface {
	Interface() (interface{}, error)
	This() *RdbResult
}

type RdbHashResultInterface interface {
	RdbResultInterface
	ScanMapByMergeArgs(interface{}) error
	ScanMap(interface{}) error
}
type RdbResult struct {
	args []interface{} // 命令参数
	data interface{}
	err  error
}

func NewErrRdbResult(format string, a ...interface{}) *RdbResult {
	return &RdbResult{
		data: nil,
		err:  errors.New(fmt.Sprintf(format, a...)),
	}
}

func (opt *RdbResult) Error() error {
	return opt.err
}

func (opt *RdbResult) Data() interface{} {
	return opt.data
}
func (opt *RdbResult) This() *RdbResult {
	return opt
}

func (opt *RdbResult) SetResult(data interface{}, err error) {
	opt.data = data
	opt.err = err
}

func (opt *RdbResult) Bytes() ([]byte, error) {
	return redis.Bytes(opt.data, opt.err)
}

func (opt *RdbResult) ByteSlices() ([][]byte, error) {
	return redis.ByteSlices(opt.data, opt.err)
}

func (opt *RdbResult) String() (string, error) {
	return redis.String(opt.data, opt.err)
}

func (opt *RdbResult) StringSlices() ([]string, error) {
	return redis.Strings(opt.data, opt.err)
}

func (opt *RdbResult) Int() (int, error) {
	return redis.Int(opt.data, opt.err)
}

func (opt *RdbResult) IntSlices() ([]int, error) {
	return redis.Ints(opt.data, opt.err)
}

func (opt *RdbResult) Float64() (float64, error) {
	return redis.Float64(opt.data, opt.err)
}

func (opt *RdbResult) Float64Slices() ([]float64, error) {
	return redis.Float64s(opt.data, opt.err)
}

func (opt *RdbResult) Bool() (bool, error) {
	return redis.Bool(opt.data, opt.err)
}

func (opt *RdbResult) IntMap() (map[string]int, error) {
	return redis.IntMap(opt.data, opt.err)
}

func (opt *RdbResult) StringMap() (map[string]string, error) {
	return redis.StringMap(opt.data, opt.err)
}

func (opt *RdbResult) Interface() (interface{}, error) {
	return opt.data, opt.err
}

func (opt *RdbResult) Values() ([]interface{}, error) {
	return redis.Values(opt.data, opt.err)
}

func (opt *RdbResult) Struct(s interface{}) error {
	bytes, err := opt.Bytes()
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, s)
	if err != nil {
		return err
	}
	return nil
}

var errScanMapValue = errors.New("rediuse.ScanMap: value must be non-nil pointer to a map")

// 将 hgetall 的结果转化为 map[T]T
func (opt *RdbResult) ScanMap(s interface{}) error {
	src, err := opt.Values()
	if err != nil {
		return err
	}
	return opt.scanMap(s, src)
}

// 将 hmget 的结果转化为 map[T]T
func (opt *RdbResult) ScanMapByMergeArgs(s interface{}) error {
	src, err := opt.Values()
	if err != nil {
		return err
	}
	newSrc := make([]interface{}, len(src)*2)
	for i := 0; i < len(src); i++ {
		newSrc[2*i] = opt.args[i+1]
		newSrc[2*i+1] = src[i]
	}
	return opt.scanMap(s, newSrc)
}

func (opt *RdbResult) scanMap(s interface{}, src []interface{}) error {

	var err error
	defer func() {
		_err := recover()
		if _err != nil {
			fmt.Printf("scanMap error:%s\n", _err)
		}
	}()

	d := reflect.ValueOf(s)
	if d.Kind() != reflect.Ptr || d.IsNil() {
		return errScanMapValue
	}
	d = d.Elem()
	if d.Kind() != reflect.Map {
		return errScanMapValue
	}
	t := d.Type()

	if len(src)%2 != 0 {
		return errors.New("rediuse.ScanMap:  number of mgethash'result not a multiple of 2")
	}

	// key处理函数，map.key的类型固定，因此循环中的处理过程都是确定的，这里应该写成多个命名函数的，但是偷懒了一点
	var ktHandle func(key string, kv *reflect.Value) error
	kt := t.Key()
	switch {
	case kt.Kind() == reflect.String:
		ktHandle = func(key string, kv *reflect.Value) error {
			*kv = reflect.ValueOf(key).Convert(kt)
			return nil
		}
	default:
		switch kt.Kind() {
		case reflect.Float32, reflect.Float64:
			ktHandle = func(key string, kv *reflect.Value) error {
				s := string(key)
				n, err := strconv.ParseFloat(s, 64)
				if err != nil || reflect.Zero(kt).OverflowFloat(n) {
					errMsg := ""
					if err != nil {
						errMsg = err.Error()
					}
					return errors.New("number " + s + " convert failed. Unmarshal Type(float) err:" + errMsg)
				}
				*kv = reflect.ValueOf(n).Convert(kt)
				return nil
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:

			ktHandle = func(key string, kv *reflect.Value) error {
				s := string(key)
				n, err := strconv.ParseInt(s, 10, 64)
				if err != nil || reflect.Zero(kt).OverflowInt(n) {
					errMsg := ""
					if err != nil {
						errMsg = err.Error()
					}
					return errors.New("number " + s + " convert failed. Unmarshal Type(int) err:" + errMsg)
				}
				*kv = reflect.ValueOf(n).Convert(kt)
				return nil
			}

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:

			ktHandle = func(key string, kv *reflect.Value) error {
				s := string(key)
				n, err := strconv.ParseUint(s, 10, 64)
				if err != nil || reflect.Zero(kt).OverflowUint(n) {
					errMsg := ""
					if err != nil {
						errMsg = err.Error()
					}
					return errors.New("number " + s + " convert failed. Unmarshal Type(uint) err:" + errMsg)
				}
				*kv = reflect.ValueOf(n).Convert(kt)
				return nil
			}

		default:

			errMsg := ""
			if err != nil {
				errMsg = err.Error()
			}
			return fmt.Errorf("convert failed. Unmarshal Type(%t) err:"+errMsg, kt.Kind())
		}
	}
	// value处理函数，map.value的类型固定，因此循环中的处理过程都是确定的，这里应该写成多个命名函数的，但是偷懒了一点
	var vtHandle func(data []byte, vv *reflect.Value) error
	elemType := t.Elem()
	switch elemTypekind := elemType.Kind(); elemTypekind {
	case reflect.Ptr, reflect.Struct, reflect.Slice, reflect.Map:
		elemTypekind = elemType.Elem().Kind()
		switch elemTypekind {
		case reflect.Struct, reflect.Slice, reflect.Map:
			vtHandle = func(data []byte, elemv *reflect.Value) error {
				err = json.Unmarshal(data, (*elemv).Interface())
				return err
			}
			break
		case reflect.String:
			vtHandle = func(key []byte, kv *reflect.Value) error {
				*kv = reflect.ValueOf(key).Convert(kt)
				return nil
			}
		case reflect.Float32, reflect.Float64:
			vtHandle = func(key []byte, kv *reflect.Value) error {
				s := string(key)
				n, err := strconv.ParseFloat(s, 64)
				if err != nil || reflect.Zero(kt).OverflowFloat(n) {
					errMsg := ""
					if err != nil {
						errMsg = err.Error()
					}
					return errors.New("number " + s + " convert failed. Unmarshal Type(float) err:" + errMsg)
				}
				*kv = reflect.ValueOf(n).Convert(kt)
				return nil
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:

			vtHandle = func(key []byte, kv *reflect.Value) error {
				s := string(key)
				n, err := strconv.ParseInt(s, 10, 64)
				if err != nil || reflect.Zero(kt).OverflowInt(n) {
					errMsg := ""
					if err != nil {
						errMsg = err.Error()
					}
					return errors.New("number " + s + " convert failed. Unmarshal Type(int) err:" + errMsg)
				}
				*kv = reflect.ValueOf(n).Convert(kt)
				return nil
			}
		case reflect.Bool:
			vtHandle = func(key []byte, kv *reflect.Value) error {
				s := string(key)
				n, err := strconv.ParseInt(s, 10, 64)
				if err != nil || reflect.Zero(kt).OverflowInt(n) {
					errMsg := ""
					if err != nil {
						errMsg = err.Error()
					}
					return errors.New("number " + s + " convert failed. Unmarshal Type(bool) err:" + errMsg)
				}
				boolV := false
				if n > 0 {
					boolV = true
				}
				*kv = reflect.ValueOf(boolV).Convert(kt)
				return nil
			}

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:

			vtHandle = func(key []byte, kv *reflect.Value) error {
				s := string(key)
				n, err := strconv.ParseUint(s, 10, 64)
				if err != nil || reflect.Zero(kt).OverflowUint(n) {
					errMsg := ""
					if err != nil {
						errMsg = err.Error()
					}
					return errors.New("number " + s + " convert failed. Unmarshal Type(uint) err:" + errMsg)
				}
				*kv = reflect.ValueOf(n).Convert(kt)
				return nil
			}

		default:
			return fmt.Errorf("can't handle type:%s\n", elemType.Kind())
		}
	default:
		return fmt.Errorf("can't handle type:%s\n", elemType.Kind())
	}

	for i := 0; i < len(src); i += 2 {

		// 处理key
		name, ok := src[i].([]byte)
		if !ok {
			name = []byte(getString(src[i]))
			//return fmt.Errorf("rediuse.ScanMap: key %d not a bulk string value", i)
		}
		key := name

		// 代码片段取自 /encoding/json/decode.go:661 func (d *decodeState) object(v reflect.Value) error
		var kv reflect.Value
		err = ktHandle(string(key), &kv)
		if err != nil {
			return err
		}

		// 处理value
		valueData, ok := src[i+1].([]byte)
		if !ok {
			if src[i+1] == nil || reflect.ValueOf(src[i+1]).IsNil() {
				continue
			}
			return fmt.Errorf("redigo.ScanStruct: key %d not a bulk string value", i)
		}
		data := valueData
		var subv reflect.Value
		var elemv reflect.Value

		if !subv.IsValid() {
			elemv = reflect.New(elemType)
			subv = elemv.Elem()
		} else {
			elemv = reflect.Zero(elemType)
			subv.Set(elemv)
		}
		if err = vtHandle(data, &elemv); err != nil {
			return err
		}

		if kv.IsValid() {
			d.SetMapIndex(kv, subv)
		}
	}
	return nil
}
