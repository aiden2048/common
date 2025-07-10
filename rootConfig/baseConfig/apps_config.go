package baseConfig

import (
	"context"

	"github.com/aiden2048/pkg/frame"
	"github.com/aiden2048/pkg/frame/logs"
	"github.com/aiden2048/pkg/public/configCache"
	"github.com/aiden2048/pkg/utils"
	jsoniter "github.com/json-iterator/go"

	"github.com/aiden2048/common/model/configModel"
	"github.com/aiden2048/common/rootConfig/base"
)

const (
	defaultRatio    = 10000 // 默认比例
	defaultCurrency = "USD"
)

type AppsConfig struct {
	base.Config
	Md5Sums map[int32]string
	allApps map[int32]*configModel.CApps
}

func InstAppsConfig(install ...interface{}) *AppsConfig {
	if len(install) > 0 && install[0].(bool) {
		inst := configCache.SetInstance("AppsConfig", &AppsConfig{
			Md5Sums: map[int32]string{},
		})

		return inst.(*AppsConfig)
	}
	inst := configCache.GetInstance("AppsConfig")
	return inst.(*AppsConfig)
}
func (c *AppsConfig) ListenTables() []string {
	return []string{"c_apps"}
}

func (c *AppsConfig) GetPower() int {
	return base.PowerAppConfig
}

func (c *AppsConfig) LoadConfig() (err error) {
	logs.PrintBill("root_config", "Begin================AppsConfig->LoadConfig=============")
	if err = c.LoadAppsConfig(); err == nil {
	}

	logs.PrintBill("root_config", "End================AppsConfig->LoadConfig=============", err, "\n")
	return

}
func (c *AppsConfig) LoadAppsConfig() error {
	mod := configModel.NewCAppsModel(int64(frame.GetPlatformId()))
	res, err := mod.FindAll(context.TODO(), nil)
	logs.PrintBill("root_config", "LoadAppsConfig", res, err)
	if err != nil && err != configModel.ErrNotFound {
		return err
	}
	if len(res) == 0 {
		app := &configModel.CApps{
			AppId:      frame.GetPlatformId()*100 + 1,
			Name:       "default",
			Country:    "BR",
			PlatId:     frame.GetPlatformId(),
			Lang:       []int32{1, 11},
			Status:     1,
			TimeZone:   0,
			PointType:  5,
			PointRatio: 10000,
			Currency:   "BRL",
			Code:       "DEFAULT",
		}
		res = append(res, app)
		if err := mod.InsertOne(context.Background(), app); err != nil {
			logs.Errorf("LoadAppsConfig err:%+v", err)
		}
	}
	allApps := make(map[int32]*configModel.CApps, len(res))
	for _, val := range res {
		allApps[val.AppId] = val
		dataStr, _ := jsoniter.MarshalToString(val)
		sum := utils.Md5(dataStr)
		if base.NeedOnlinePush() && c.Md5Sums[val.AppId] != sum {
			c.Md5Sums[val.AppId] = sum
			base.PushToOnline(val.AppId, c.ListenTables())
		}
	}
	c.allApps = allApps

	return nil
}

func (c *AppsConfig) GetAppInfo(appid int32) *configModel.CApps {
	app, ok := c.allApps[appid]
	if ok {
		return app
	}

	return nil // &configModel.CApps{}

}

// GetAppPointRatio app配置的货币兑换比例
func (c *AppsConfig) GetAppPointRatio(appid int32) int64 {
	app, ok := c.allApps[appid]
	if !ok {
		return defaultRatio
	}
	if app != nil && app.PointRatio > 0 {
		return int64(app.PointRatio)
	}
	return defaultRatio
}

// GetAppCurrency app配置的货币代码
func (c *AppsConfig) GetAppCurrency(appid int32) string {
	app, ok := c.allApps[appid]
	if !ok {
		return defaultCurrency
	}
	if app != nil && app.Currency != "" {
		return app.Currency
	}
	return defaultCurrency
}
func (c *AppsConfig) GetAppCode(appid int32) string {
	app, ok := c.allApps[appid]
	if !ok {
		return ""
	}
	if app != nil {
		return app.Code
	}
	return ""

}
func (c *AppsConfig) GetAllApps() map[int32]*configModel.CApps {
	return c.allApps
}

// GetAppPointType 获取app钱包币种id
func (c *AppsConfig) GetAppPointType(appid int32) (int32, bool) {
	appInfo, ok := c.allApps[appid]
	if !ok {
		return 0, ok
	}
	return appInfo.PointType, ok
}
