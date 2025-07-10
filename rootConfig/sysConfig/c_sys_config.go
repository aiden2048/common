package sysConfig

import (
	"context"
	"fmt"
	"sync"

	"github.com/aiden2048/common/rootConfig/base"

	"github.com/aiden2048/common/model/configModel"

	"github.com/aiden2048/pkg/public/configCache"

	"github.com/aiden2048/pkg/frame"
	"github.com/aiden2048/pkg/frame/logs"
	"github.com/aiden2048/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"

	jsoniter "github.com/json-iterator/go"
)

type SysConfig struct {
	base.Config
	allConfig *sync.Map
	Md5Sums   map[int32]string
}

func InstSysConfig(install ...interface{}) *SysConfig {
	if len(install) > 0 && install[0].(bool) {
		inst := configCache.SetInstance("SysConfig", &SysConfig{
			allConfig: &sync.Map{},
			Md5Sums:   map[int32]string{},
		})
		isInst = true
		//inst.SetPower(base.PowerSysConfig)
		return inst.(*SysConfig)
	}
	inst := configCache.GetInstance("SysConfig")
	return inst.(*SysConfig)
}
func (c *SysConfig) ListenTables() []string {
	return []string{"c_config"}
}
func (c *SysConfig) GetPower() int {
	return base.PowerSysConfig
}
func (c *SysConfig) LoadConfig() error {
	if err := c.loadConfig(); err != nil {
		return err
	}
	ReloadSubConfig()
	return nil
}
func (c *SysConfig) loadConfig() error {
	rCnfModel := configModel.NewCConfigModel(int64(frame.GetPlatformId()))
	// 读新库
	rets, err := rCnfModel.FindAll(context.TODO(), bson.M{})
	if err != nil && err != configModel.ErrNotFound {
		logs.Errorf("loadConfig rCnfModel.FindAll err:%+v", err)
		return nil
	}
	allConfig := &sync.Map{}
	datas := map[int32][]any{}
	for _, ret := range rets {
		if ret.Value == "" {
			continue
		}
		key := fmt.Sprintf("%d_%d_%s", ret.PlatId, ret.AppId, ret.Key)
		allConfig.Store(key, ret.Value)
		datas[ret.AppId] = append(datas[ret.AppId], ret)
	}
	c.allConfig = allConfig
	for app_id, data := range datas {
		if app_id == 0 {
			continue
		}
		data = append(data, datas[0]...)
		dataStr, _ := jsoniter.MarshalToString(data)
		sum := utils.Md5(dataStr)
		if base.NeedOnlinePush() && c.Md5Sums[app_id] != sum {
			c.Md5Sums[app_id] = sum
			base.PushToOnline(app_id, c.ListenTables())
		}
	}

	return nil
}

// 读取三次，第一次读取自己app_id的数据，如果没有话再读取 app_id = 0
func (c *SysConfig) GetConfig(app_id int32, key string) string {
	mapKey := fmt.Sprintf("%d_%d_%s", frame.GetPlatformId(), app_id, key)
	val, ok := c.allConfig.Load(mapKey)
	if ok {
		if cfg, ok := val.(string); ok {
			return cfg
		}
	}
	mapKey = fmt.Sprintf("%d_%d_%s", frame.GetPlatformId(), 0, key)
	val, ok = c.allConfig.Load(mapKey)
	if ok {
		if cfg, ok := val.(string); ok {
			return cfg
		}
	}
	mapKey = fmt.Sprintf("%d_%d_%s", 0, 0, key)
	val, ok = c.allConfig.Load(mapKey)
	if ok {
		if cfg, ok := val.(string); ok {
			return cfg
		}
	}
	logs.Importantf("GetConfig failed!!! err: Not found, appId:%d, key:%s", app_id, key)
	return ""
}

func (c *SysConfig) FindAndUnmarshalConfig(appId int32, key string, val interface{}) (bool, error) {
	cfgStr := c.GetConfig(appId, key)
	if cfgStr == "" {
		logs.Debugf("FindAndUnmarshalConfig GetConfig cfgStr == nil, appId:%d, key:%s", appId, key)
		return false, nil
	}

	err := jsoniter.Unmarshal([]byte(cfgStr), val)
	if err != nil {
		logs.Errorf("FindAndUnmarshalConfig err:%+v, appId:%d, key:%s", err, appId, key)
		return false, err
	}

	return true, nil
}
