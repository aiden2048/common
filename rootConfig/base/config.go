package base

import (
	"github.com/aiden2048/pkg/public/configCache"
	"github.com/aiden2048/pkg/utils"
)

const (
	//加载排序, 越大越早执行
	PowerBeginConfig     = 99999999  //最后一个执行的
	PowerEndConfig       = -99999999 //最后一个执行的
	PowerAppConfig   int = 1000 - iota
	PowerSysConfig
	PowerOriGameConfig
	PowerExtGameRouteConfig
	PowerExtGameConfig
	PowerExtAppGameConfig
	PowerPackConfig
)

type Config struct {
	confgKey string
}

func (c *Config) GetPower() int { // 默认间隔时间60s
	return 0
}

// func (c *Config) SetPower(p int) { // 默认间隔时间60s
//
//		c.power = p
//	}
func (c *Config) Intervals() int64 { // 默认间隔时间60s
	return 60
}
func (c *Config) SetConfigKey(k string) {
	c.confgKey = k
}
func (c *Config) GetConfigKey() string {
	return c.confgKey
}
func (c *Config) ListenTables() []string {
	return []string{}
}
func (c *Config) LoadConfig() (err error) {
	return nil
}
func (c *Config) ReloadConfig() (configCache.ConfigInterface, error) {
	var inst configCache.ConfigInterface
	if c.confgKey != "" {
		inst = configCache.GetInstance(c.confgKey)
	}
	if inst == nil {
		panic(utils.GetCallFile(1) + "---" + c.confgKey + "---没有实现 ReloadConfig")
	}
	return inst, inst.LoadConfig()
}
