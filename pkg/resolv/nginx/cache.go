package nginx

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

type Caches map[string]cache

func NewCaches() Caches {
	return Caches{}
}

func (cs *Caches) SetCache(config *Config, file ...interface{}) error {
	cache, err := newCache(config, file...)
	if err != nil {
		return err
	}
	oldCache, ok := (*cs)[config.Value]
	if ok && oldCache.hash == cache.hash {
		return IsInCaches
	}
	(*cs)[config.Value] = cache
	return nil
}

func (cs Caches) IsCached(path string) bool {
	_, ok := cs[path]
	return ok
}

func (cs Caches) CheckHash(path string) (bool, error) {
	if cs.IsCached(path) {
		hash, hashErr := getHash(path)
		return hash == cs[path].hash, hashErr
	}
	return false, IsNotInCaches
}

func (cs Caches) GetConfig(path string) (*Config, error) {
	if cache, ok := cs[path]; ok {
		return cache.config, nil
	}
	return nil, fmt.Errorf("the Config(path: %s) cache object does not exist", path)
}

type cache struct {
	config *Config
	hash   string
}

func newCache(config *Config, file ...interface{}) (cache, error) {
	hash, err := getHash(config.Value, file...)
	return cache{
		config: config,
		hash:   hash,
	}, err
}

// getHash, 计算文件hash值函数
// 参数:
//
//	config: Config对象指针
//
// 返回值:
//
//	文件哈希基准值
//	错误
func getHash(path string, file ...interface{}) (hash string, err error) {
	hash256 := sha256.New()
	// fmt.Println("------------getHash func------------")
	// 结束操作前捕捉panic
	defer func() {
		// 捕获panic
		if r := recover(); r != nil {
			// err = fmt.Errorf("%s", r)
			fmt.Println(r)
		}
		// fmt.Println("----------getHash func end----------")
	}()
	var f *os.File
	if file == nil {
		// 测试
		// fmt.Println("no param, path:", path)
		// 读取文件
		f, err = os.Open(path)
		if err != nil {
			return "", err
		}
		defer f.Close()
	} else {
		data := file[0]
		switch data.(type) {
		case string:
			// 测试
			// fmt.Println("hash:", data.(string), "path:", path)
			return data.(string), nil
		case *os.File:
			// 测试
			// fmt.Println("fd, path:", path)
			f = data.(*os.File)
			_, _ = f.Seek(0, 0)
		case []byte:
			// 测试
			// fmt.Println("bytes, path:", path)
			r := bytes.NewReader(data.([]byte))
			_, hashCPErr := io.Copy(hash256, r)
			if hashCPErr != nil {
				return "", hashCPErr
			}

			return hex.EncodeToString(hash256.Sum(nil)), nil
		default:
			return "", fmt.Errorf("wrong type or wrong value for input to func pkg/resolv/nginx.getHash")
		}
	}

	// 计算文件数据哈希值
	_, hashCPErr := io.Copy(hash256, f)
	if hashCPErr != nil {
		return "", hashCPErr
	}
	return hex.EncodeToString(hash256.Sum(nil)), nil
}
