package baseConfig

import (
	"context"
	"fmt"
	"sync"

	"github.com/aiden2048/common/rootConfig/base"

	"github.com/aiden2048/common/model/configModel"
	"github.com/aiden2048/pkg/frame"
	"github.com/aiden2048/pkg/frame/logs"
	"github.com/aiden2048/pkg/public/configCache"
	"go.mongodb.org/mongo-driver/bson"
)

type ServiceOptCodeConfig struct {
	base.Config
	allConfig *sync.Map
}

func InstServiceOptCodeConfig(install ...interface{}) *ServiceOptCodeConfig {

	if len(install) > 0 && install[0].(bool) {
		inst := configCache.SetInstance("ServiceOptCodeConfig", &ServiceOptCodeConfig{
			allConfig: &sync.Map{},
		})

		return inst.(*ServiceOptCodeConfig)
	}
	inst := configCache.GetInstance("ServiceOptCodeConfig")
	return inst.(*ServiceOptCodeConfig)
}
func (c *ServiceOptCodeConfig) ListenTables() []string {
	return []string{"c_service_optcode"}
}
func (c *ServiceOptCodeConfig) GetPower() int {
	return base.PowerSysConfig
}
func (c *ServiceOptCodeConfig) LoadConfig() (err error) {
	logs.PrintBill("root_config", "Begin================ServiceOptCodeConfig->LoadConfig=============")

	c.loadAppConfig()

	logs.PrintBill("root_config", "End================ServiceOptCodeConfig->LoadConfig=============", err, "\n")
	return

}

func (c *ServiceOptCodeConfig) loadAppConfig() {
	rCnfModel := configModel.NewCServiceOptcodeModel(int64(frame.GetPlatformId()))
	rets, _ := rCnfModel.FindAll(context.TODO(), bson.M{})
	if len(rets) == 0 {
		logs.PrintBill("root_config", "ServiceOptCodeConfig->LoadConfig============= retd:%+v", rets, "\n")
		return
	}
	for _, ret := range rets {
		c.allConfig.Store(fmt.Sprintf("%d.%d", ret.AppId, ret.Code), ret)
	}
}

func (c *ServiceOptCodeConfig) GetOptCode(app_id int32, opt_code int) *configModel.CServiceOptcode {
	key := fmt.Sprintf("%d.%d", app_id, opt_code)
	val, ok := c.allConfig.Load(key)
	if ok {
		ret, ok := val.(*configModel.CServiceOptcode)
		if ok {
			return ret
		}
	}
	key = fmt.Sprintf("%d.%d", 0, opt_code)
	val, ok = c.allConfig.Load(key)
	if ok {
		ret, ok := val.(*configModel.CServiceOptcode)
		if ok {
			return ret
		}
	}
	return nil
}
