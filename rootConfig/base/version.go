package base

import (
	"fmt"
	"sync"
	"time"

	"github.com/aiden2048/pkg/frame/logs"
	"github.com/aiden2048/pkg/public/redisDeal"
	"github.com/aiden2048/pkg/public/redisKeys"
)

var needOnlinePush bool
var version *sync.Map
var appVersionChange *sync.Map

func SetOnlinePush() {
	if !needOnlinePush {

		needOnlinePush = true
		version = &sync.Map{}
		logs.Importantf("baseConfigVersion:SetOnlinePush NeedOnlinePush:%v", needOnlinePush)
		InstVersionConfig(true)
	}
}
func NeedOnlinePush() bool {
	return needOnlinePush
}
func ClearAppVersion() {
	logs.Print("baseConfigVersion:开始重新计算每个app配置的版本号")
	appVersionChange = &sync.Map{}
}
func GetVerSion(app_id int32) string {
	if version == nil {
		version = &sync.Map{}
	}
	v, ok := version.Load(app_id)
	if !ok {
		key := GetConfigVersionKey(app_id)
		ver := redisDeal.RedisDoGetStr(key)
		if ver == "" {
			ver = newVersion()
			// 写redis
			redisDeal.RedisDoSet(key, ver, redisDeal.InfoTtlTwoMonth)
		}
		logs.Print("baseConfigVersion:保存版本号", app_id, ver)
		version.Store(app_id, ver)
		return ver
	}
	return fmt.Sprintf("%v", v)
}
func PushToOnline(app_id int32, tables []string) {
	if app_id == 0 || !needOnlinePush {
		return
	}
	logs.PrintBill("AppVersionChangePushToOnline", app_id, tables)
	if appVersionChange == nil {
		appVersionChange = &sync.Map{}
	}
	appVersionChange.Store(app_id, true)
}

func PushMsgToOnline(app_id int32) {

	ver := newVersion()
	logs.Importantf("baseConfigVersion:PushToOnline app_id:%d,version:%s", app_id, ver)
	//onlineapi.PushMsgToAll(commonConst.ONLINE_EVENT_CONFIG_CHANGE, ver, app_id, nil)

	key := GetConfigVersionKey(app_id)
	redisDeal.RedisDoSet(key, ver, redisDeal.InfoTtlTwoMonth)

}
func newVersion() string {
	ver := fmt.Sprintf("%d", time.Now().Unix())
	return ver
}
func ReportAllAppVersion() bool {
	if appVersionChange == nil {
		return false
	}
	logs.Print("baseConfigVersion:ReportAllAppVersion")
	count := 0
	appVersionChange.Range(func(key, value any) bool {
		appid := key.(int32)
		PushMsgToOnline(appid)
		count++
		return true
	})
	return count > 0
}

func ReloadAllAppVersion() {
	logs.Print("baseConfigVersion:ReloadAllAppVersion")
	version = &sync.Map{}
}
func GetConfigVersionKey(app_id int32) *redisKeys.RedisKeys {
	tmp := &redisKeys.RedisKeys{}
	tmp.Name = redisKeys.REDIS_INDEX_COMMON
	tmp.Key = fmt.Sprintf("rootconfig.version.%d", app_id)
	return tmp
}
