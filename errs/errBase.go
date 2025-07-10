package errs

import (
	"github.com/aiden2048/pkg/frame/logs"
	"github.com/aiden2048/pkg/public/errorMsg"
)

func NewError(errCodes []int32, ret int32, msgfmt string, only_back_end bool, p ...interface{}) *errorMsg.ErrRsp { // 有fmt參數時使用
	if len(errCodes) < 2 {
		logs.Errorf("错误码超出范围 ret:%d,msgfmt:%s", ret, msgfmt)
		return errorMsg.NewError(ret, msgfmt, p...)
	}
	if only_back_end {
		return errorMsg.NewError(ret+errCodes[0], msgfmt, p...)
	}
	if ret+errCodes[0] >= errCodes[1] {
		logs.Errorf("错误码超出范围 ret:%d,msgfmt:%s,errCodes:%+v", ret, msgfmt, errCodes)
	}
	return errorMsg.NewConstError(ret+errCodes[0], msgfmt, p...)
}
