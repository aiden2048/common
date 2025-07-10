package {{.pkgName}}

import (
	"application/services/activity/repeatActivityService/baseActivity"
	"application/services/activity/types"
	"context"
	"fmt"

	"github.com/aiden2048/common/errs/commErr"
	"github.com/aiden2048/common/model/configModel"
	"github.com/aiden2048/common/public/events"
	"github.com/aiden2048/common/serviceApi/jobSchedApi"
	"github.com/aiden2048/common/serviceApi/lockerSvrApi"
	"github.com/aiden2048/pkg/public/errorMsg"
	"github.com/aiden2048/pkg/public/redisKeys"
	"github.com/aiden2048/pkg/qgframe"
	"github.com/aiden2048/pkg/qgframe/logs"
	jsoniter "github.com/json-iterator/go"
)

// 活动
type Activity struct {
	baseActivity.Activity
	selfConfig SelfConfig
}

func NewService(ctx context.Context, sess *qgframe.Session, config *configModel.CRepeatActivityConfig) baseActivity.ActivityService {
	cfg := SelfConfig{}
	if config.Config == "" {
		return nil
	}
	err := jsoniter.UnmarshalFromString(config.Config, &cfg)
	if err != nil {
		logs.Errorf("subscribe json解析失败,err:%+v config:%s", err, config.Config)
		return nil
	}
	return &Activity{
		Activity: baseActivity.Activity{
			Sess:   sess,
			Ctx:    ctx,
			Config: config,
		},
		selfConfig: cfg,
	}
}

func (a *Activity) GetConfig(s string) (any, *errorMsg.ErrRsp) {
	if s == "" {
		return a.selfConfig, nil
	}

	cfg := SelfConfig{}
	err := jsoniter.UnmarshalFromString(s, &cfg)
	if err != nil {
		logs.Errorf("err:%v", err)
		return nil, commErr.QueryDataErr.Return(err)
	}
	return cfg, nil
}

func (a *Activity) DoActive(req *types.ActivityReq) (rsp *types.ActiveRsp, err *errorMsg.ErrRsp) {
	rsp = &types.ActiveRsp{}
	// 加锁 防止重复领取

	lockkey := redisKeys.GenLockRedisKey(uint64(a.Config.AppId), req.UserId, fmt.Sprintf("DoActive.%d", a.Config.AType))
	lock := lockerSvrApi.QLock(lockkey.Key, 10)
	if lock == nil {
		return nil, commErr.QueryTooQuick
	}
	defer lock.Unlock()

	return
}

func (a *Activity) GetDetails() (rsp *types.GetDetailsRsp, err *errorMsg.ErrRsp) {
	rsp = &types.GetDetailsRsp{}
	rsp.Details = SelfDetails{}
	return
}

// 注册事件
func (a *Activity) RegEvent(event events.RegisterUserEvent) {

	return
}

// 充值成功事件
func (a *Activity) PayEvent(event events.UserPayPointEvent) {
	return
}

// 游戏事件
func (a *Activity) GameEvent(event events.TUserGameRecordEvent) {
	return
}

func (a *Activity) DoJob(_ *jobSchedApi.ExecJobReq, _ string) *errorMsg.ErrRsp {

	return nil
}
func (a *Activity) InitJob(a_type int32) *jobSchedApi.PubJobReq {
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

func (a *Activity) PutEvent(event events.UserPutPointEvent) {

}

func (a *Activity) InviteEvent(event events.UserInviteEvent) {

}

