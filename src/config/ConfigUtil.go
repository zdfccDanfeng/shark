package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/shark/src/util"
	"path/filepath"
	"sync"
)

// 数据库配置信息
type DbConf struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
}

func ToString(conf *AppConf) map[string]DbConf {
	for k, v := range conf.Dbs {
		fmt.Println("v is :", v, "k is :", k)
	}
	return conf.Dbs
}

type AppConf struct {
	Dbs map[string]DbConf
}

var (
	cfg     *AppConf
	once    sync.Once
	cfgLock = new(sync.RWMutex)
)

// 单例模式加载
func Config() *AppConf {
	once.Do(ReloadConfig)
	cfgLock.RLock()
	defer cfgLock.RUnlock()
	return cfg
}

func ReloadConfig() {
	pwd := util.RelativePath()
	filePath, err := filepath.Abs(pwd + "/conf/app")
	if err != nil {
		panic(err)
	}
	fmt.Printf("parse toml file once. filePath: %s\n", filePath)
	config := new(AppConf)
	if _, err := toml.DecodeFile(filePath, config); err != nil {
		panic(err)
	}
	cfgLock.Lock()
	defer cfgLock.Unlock()
	cfg = config
}
