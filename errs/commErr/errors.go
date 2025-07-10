package commErr

import (
	"github.com/aiden2048/common/errs"
	"github.com/aiden2048/pkg/public/errorMsg"
)

func NewError(ret int32, msgfmt string, only_back_end bool, p ...interface{}) *errorMsg.ErrRsp { // 有fmt參數時使用
	return errs.NewError(errs.CommErrCode, ret, msgfmt, only_back_end, p...)
}

var (
	ParamErr      = NewError(1, "参数异常", false)
	QueryTooQuick = NewError(2, "请求频率太快", false)
	UpdateDataErr = NewError(3, "保存失败,请重试", false)
	QueryDataErr  = NewError(4, "读取数据失败,请重试", false)
	TokenDelErr   = NewError(5, "登录账号被踢出", false)
	RequestErr    = NewError(6, "请求失败", false)
)
var (
	QueryDataMaxOne = NewError(1, "Only one day's data can be queried!", true)
)
