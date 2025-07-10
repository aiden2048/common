package {{.pkgName}}

import (
	"application/services/activity/activityService/baseActivity"
	"application/services/activity/types"
	"context"
	"fmt"

	"github.com/aiden2048/common/errs/commErr"

	"github.com/aiden2048/common/userCache"

	"github.com/aiden2048/common/model/configModel"

	"github.com/aiden2048/common/public/events"
	"github.com/aiden2048/common/serviceApi/jobSchedApi"
	"github.com/aiden2048/common/serviceApi/lockerSvrApi"
	"github.com/aiden2048/pkg/public/errorMsg"
	"github.com/aiden2048/pkg/public/redisDeal"
	"github.com/aiden2048/pkg/public/redisKeys"
	"github.com/aiden2048/pkg/frame/logs"
	jsoniter "github.com/json-iterator/go"
)

type Activity struct {
	baseActivity.Activity
	ctx        context.Context
	config     *configModel.CActivityConfig
	selfConfig SelfConfig
}

func (l *Activity) NewActivity(ctx context.Context, config *configModel.CActivityConfig) baseActivity.ActivityService {
	cfg := SelfConfig{}
	if config.Config == "" {
		return nil
	}
	err := jsoniter.UnmarshalFromString(config.Config, &cfg)
	if err != nil {
		logs.Errorf("json解析失败,err:%+v config:%s", err, config.Config)
		return nil
	}
	l.config = config
	return &Activity{
		ctx:        ctx,
		config:     config,
		selfConfig: l.makeNewConfig(cfg),
	}
}

func (l *Activity) GetNewConfig(cfgStr string, version int32) (rsp any, err error) {
	if cfgStr == "" {
		cfgStr = l.config.Config
	}
	cfg := SelfConfig{}
	err = jsoniter.UnmarshalFromString(cfgStr, &cfg)
	if err != nil {
		return nil, err
	}
	return l.makeNewConfig(cfg), nil
}
func (l *Activity) makeNewConfig(cfg SelfConfig) SelfConfig {

	return cfg
}
func (l *Activity) DoActive(userInfo *userCache.UserSimpleInfo, req *types.ActivityReq) (rsp *types.ActiveRsp, err *errorMsg.ErrRsp) {
	rsp = &types.ActiveRsp{}
	// 加锁 防止重复领取
	lockkey := redisKeys.GenLockRedisKey(uint64(l.config.AppId), userInfo.UserId, fmt.Sprintf("DoActive.%d", l.config.AType))
	lock := lockerSvrApi.QLock(lockkey.Key, 10)
	if lock == nil {
		return nil, commErr.QueryTooQuick
	}
	defer lock.Unlock()

	return
}

func (l *Activity) GetDetails(userInfo *userCache.UserSimpleInfo) (rsp *types.ActivityDetails, err *errorMsg.ErrRsp) {
	rsp = &types.ActivityDetails{}
	rsp.Details = SelfDetails{}
	return
}

// 注册事件
func (l *Activity) RegEvent(userInfo *userCache.UserSimpleInfo, event events.RegisterUserEvent) {

	return
}

// 充值成功事件
func (l *Activity) PayEvent(userInfo *userCache.UserSimpleInfo, event events.UserPayPointEvent) {
	return
}

// 游戏事件
func (l *Activity) GameEvent(userInfo *userCache.UserSimpleInfo, event events.TUserGameRecordEvent) {
	return
}

func (l *Activity) DoJob(_ *jobSchedApi.ExecJobReq, _ string) *errorMsg.ErrRsp {

	return nil
}
func (l *Activity) InitJob(a_type int32) *jobSchedApi.PubJobReq {
	// TimeCronExpr := "0 10 2 * * * *"
	// return &jobSchedApi.PubJobReq{
	// 	Name:                         fmt.Sprintf("activity.job.%d", a_type),
	// 	AllowReRunCntCazFailPerSched: 600,
	// 	Mode:                         commonConst.SchedJobModeTimeCron,
	// 	TimeCronExpr:                 TimeCronExpr,
	// 	UniqBillNo:                   fmt.Sprintf("activity.job.%d", a_type),
	// }
	return nil
}

func (l *Activity) PutEvent(userInfo *userCache.UserSimpleInfo, event events.UserPutPointEvent) {

}

func (l *Activity) InviteEvent(userInfo *userCache.UserSimpleInfo, event events.UserInviteEvent) {

}
