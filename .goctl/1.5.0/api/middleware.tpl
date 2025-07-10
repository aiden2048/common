package middleware

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/aiden2048/common/serviceApi/connApi"
	"github.com/aiden2048/pkg/public/httpHandle"
	"github.com/aiden2048/pkg/frame"
	"github.com/aiden2048/pkg/frame/logs"
	"github.com/aiden2048/pkg/utils"
	"github.com/aiden2048/pkg/utils/baselib"
	"github.com/valyala/fasthttp"
)

var NoLoginFunc = []string{"login"}

func HttpMiddleware(ctx *fasthttp.RequestCtx, funcName string, next func(ctx context.Context, r *frame.NatsMsg)) {
	msg := &frame.NatsMsg{}

	body := ctx.Request.Body()
	msg.MsgData = body

	herder := map[string][]string{}
	ctx.Request.Header.VisitAll(func(key, value []byte) {
		herder[string(key)] = []string{string(value)}
	})
	// 校验token
	err, sess := CheckLogin(ctx, funcName)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		httpHandle.SendResponse(ctx, 1, err.Error(), nil)
		return
	}
	msg.Sess = *sess
	lang, _ := strconv.Atoi(string(ctx.QueryArgs().Peek("lang")))
	msg.Sess.AppType = 1
	msg.Sess.Lang = int32(lang)
	msg.Sess.Os = 1
	msg.Sess.Addr = baselib.GetAddresFromHttpHeader(herder, ctx.RemoteAddr().String())
	msg.Sess.MutableLoginInfo().RemoteAddr = msg.Sess.Addr
	msg.Sess.Host = string(ctx.Host())
	gpid := int64(0)
	if msg.Sess.Trace == nil {
		msg.Sess.Trace, gpid = logs.CreateTraceId(&msg.Sess)
	} else {
		gpid = logs.StoreTraceId(msg.Sess.Trace)
	}
	defer func() {
		if gpid != 0 {
			logs.RemoveTraceId(gpid)
		}
	}()
	next(context.TODO(), msg)
}
func CheckLogin(ctx *fasthttp.RequestCtx, funcName string) (err error, sess *frame.Session) {
	token := string(ctx.Request.Header.Peek("token"))
	if token == "" {
		if !utils.InArray(NoLoginFunc, funcName) {
			err = errors.New("no token")
			return
		}
		appId, _ := strconv.Atoi(string(ctx.Request.Header.Peek("app_id")))
		sess = frame.NewSessionOnlyApp(int32(appId))
		return
	}
	err, sess = connApi.CheckLoginV2(token, []string{connApi.SRC_ADMIN}, time.Now(), true)
	if err != nil {
		return
	}
	// 权限校验（角色/菜单/按钮）
	if !checkPermissions(ctx, funcName, sess) {
		err = errors.New("no permissions")
		return
	}
	return nil, sess
}
func checkPermissions(ctx *fasthttp.RequestCtx, funcName string, sess *frame.Session) bool {
	if strings.HasPrefix(funcName, "NA.") {
		return true
	}
	return false
}
