package base

import (
	"github.com/aiden2048/pkg/frame"
	"github.com/aiden2048/pkg/frame/logs"
	"github.com/aiden2048/pkg/public/configCache"
	"github.com/aiden2048/pkg/public/elect"
)

// 用于生成配置版本号的虚拟配置类
type ConfigVersion struct {
	Config
	Power int
}

var elc *elect.RedisElect

func InstVersionConfig(install ...interface{}) {
	if len(install) > 0 && install[0].(bool) {
		_ = configCache.SetInstance("BeginConfigVersion", &ConfigVersion{Power: PowerBeginConfig})
		_ = configCache.SetInstance("EndConfigVersion", &ConfigVersion{Power: PowerEndConfig})
		//_ = configCache.SetInstance("ReloadConfigVersion", &ReloadConfigVersion{Power: PowerEndConfig})
		//return inst.(*ConfigVersion)
	}
	//inst := configCache.GetInstance("BeginConfigVersion")

	if len(install) > 0 && install[0].(bool) {

		//return inst.(*ConfigVersion)
	}
	//inst := configCache.GetInstance("EndConfigVersion")
	//return inst.(*ConfigVersion)
	elc, _ = elect.NewRedisElect(frame.GetServerName()+"_configVersion", frame.GetStrServerID())
	elc.Run()
	frame.ListenConfig("__ReloadAppConfigVersion__", func(bytes []byte) {
		ReloadAllAppVersion()
	})
}

func (c *ConfigVersion) ListenTables() []string {
	return []string{"*"}
}
func (c *ConfigVersion) LoadConfig() error {
	if err := c.loadConfig(); err != nil {
		return err
	}
	return nil
}
func (c *ConfigVersion) loadConfig() error {
	if c.Power > 0 {
		logs.Print("开始重载配置, 准备计算版本号")
		ClearAppVersion()
	} else {
		if elc != nil && elc.IsMaster() {
			if ReportAllAppVersion() {
				frame.NotifyReloadConfig("__ReloadAppConfigVersion__", nil)
			}

		}
	}
	return nil
}
func (c *ConfigVersion) GetPower() int {
	return c.Power
}
