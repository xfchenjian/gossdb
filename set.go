package gossdb

import (
	//	"fmt"
	"github.com/seefan/goerr"
)

//设置指定 key 的值内容
//
//  key 键值
//  val 存贮的 value 值,val只支持基本的类型，如果要支持复杂的类型，需要开启连接池的 Encoding 选项
//  ttl 可选，设置的过期时间，单位为秒
//  返回 err，可能的错误，操作成功返回 nil
func (this *Client) Set(key string, val interface{}, ttl ...int) (err error) {
	var resp []string
	if len(ttl) > 0 {
		resp, err = this.Client.Do("setx", key, this.encoding(val, false), ttl[0])
	} else {
		resp, err = this.Client.Do("set", key, this.encoding(val, false))
	}
	if err != nil {
		return goerr.NewError(err, "Set %s error", key)
	}
	if len(resp) == 2 && resp[0] == "ok" {
		return nil
	}
	return makeError(resp, key)
}

//获取指定 key 的值内容
//
//  key 键值
//  返回 一个 Value,可以方便的向其它类型转换
//  返回 一个可能的错误，操作成功返回 nil
func (this *Client) Get(key string) (Value, error) {
	resp, err := this.Client.Do("get", key)
	if err != nil {
		return "", goerr.NewError(err, "Get %s error", key)
	}
	if len(resp) == 2 && resp[0] == "ok" {
		return Value(resp[1]), nil
	}
	return "", makeError(resp, key)
}

//设置过期
//
//  key 要设置过期的 key
//  ttl 存活时间(秒)
//  返回 re，设置是否成功，如果当前 key 不存在返回 false
//  返回 err，执行的错误，操作成功返回 nil
func (this *Client) Expire(key string, ttl int) (re bool, err error) {
	resp, err := this.Do("expire", key, ttl)
	if err != nil {
		return false, goerr.NewError(err, "Expire %s error", key)
	}
	if len(resp) == 2 && resp[0] == "ok" {
		return resp[1] == "1", nil
	}
	return false, makeError(resp, key, ttl)
}

//查询指定 key 是否存在
//
//  key 要查询的 key
//  返回 re，如果当前 key 不存在返回 false
//  返回 err，执行的错误，操作成功返回 nil
func (this *Client) Exists(key string) (re bool, err error) {
	resp, err := this.Do("exists", key)
	if err != nil {
		return false, goerr.NewError(err, "Exists %s error", key)
	}

	if len(resp) == 2 && resp[0] == "ok" {
		return resp[1] == "1", nil
	}
	return false, makeError(resp, key)
}

//删除指定 key
//
//  key 要删除的 key
//  返回 err，执行的错误，操作成功返回 nil
func (this *Client) Del(key string) error {
	resp, err := this.Do("del", key)
	if err != nil {
		return goerr.NewError(err, "Del %s error", key)
	}

	//response looks like this: [ok 1]
	if len(resp) > 0 && resp[0] == "ok" {
		return nil
	}
	return makeError(resp, key)
}

//指定key的原子增长
//
//key 指定的key
// vaule 增加值
//返回 err 执行的错误，err为空，返回增加后的值
func (this *Client) Incr(key string, value int) (Value, error) {
    resp, err := this.Client.Do("incr", key, value)
    if err != nil {
        return "", goerr.NewError(err, "Incr %s error", key)
    }   
    if len(resp) == 2 && resp[0] == "ok" {
        return Value(resp[1]), nil 
    }   
    return "", makeError(resp, key)
}
