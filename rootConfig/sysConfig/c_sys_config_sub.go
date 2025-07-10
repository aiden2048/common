package sysConfig

import (
	"log"
	"sync"

	"github.com/aiden2048/pkg/frame/logs"
)

var subConfigs = &sync.Map{}
var isInst = false

type SysSubInterface interface {
	LoadConfig()
}

// SysSubConfig 子类基础类
type SysSubConfig[T any] struct {
	Cache   *sync.Map
	SysKey  string
	InitCfg func(app_id int32, cfg *T)
}

// GetCfg 获取配置
func (c *SysSubConfig[T]) GetCfg(app_id int32) *T {
	t, ok := c.Cache.Load(app_id)
	if ok && t != nil {
		if cfg, ok := t.(*T); ok {
			return cfg
		}
	}
	var cfg T
	InstSysConfig().FindAndUnmarshalConfig(app_id, c.SysKey, &cfg)
	logs.Bill("c_config", "app_id:%d, c.SysKey:%s, cfg:%+v", app_id, c.SysKey, cfg)
	c.Cache.Store(app_id, &cfg)
	if c.InitCfg != nil {
		c.InitCfg(app_id, &cfg)
		logs.Bill("c_config", "InitCfg, app_id:%d, c.SysKey:%s, cfg:%+v", app_id, c.SysKey, cfg)
	}
	return &cfg
}
func (c *SysSubConfig[T]) LoadConfig() {
	c.Cache = &sync.Map{}
}

func RegistSubConfig(key string, subConf SysSubInterface) {
	logs.Debugf("RegistSubConfig config key:%+v", key)
	if !isInst {
		isInst = true
		InstSysConfig(true)
	}
	subConfigs.Store(key, subConf)
}
func GetSubConfig(key string) any {
	if !isInst {
		InstSysConfig(true)
		isInst = true
	}
	val, ok := subConfigs.Load(key)
	if !ok {
		logs.Errorf("GetSubConfig failed!!! not find config instance key:%s", key)
		log.Panicf("not find config instance: %s", key)
	}
	return val
}
func ReloadSubConfig() {
	logs.Debugf("ReloadSubConfig")
	subConfigs.Range(func(key, value any) bool {
		if subConf, ok := value.(SysSubInterface); ok {
			subConf.LoadConfig()
		}
		return true
	})
}
