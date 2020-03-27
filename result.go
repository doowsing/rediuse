package rediuse

import (
	"encoding/json"
	"github.com/gomodule/redigo/redis"
)

type RdbResult struct {
	data interface{}
	err  error
}

func NewErrRdbResult(err error) *RdbResult {
	return &RdbResult{
		data: nil,
		err:  err,
	}
}

func (opt *RdbResult) Error() error {
	return opt.err
}

func (opt *RdbResult) Data() interface{} {
	return opt.data
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
