package {{.PkgName}}

import (
	"github.com/aiden2048/pkg/public/errorMsg"
	"github.com/aiden2048/pkg/public/natsHandle"
	{{if .FromTop}}{{else}}"github.com/aiden2048/pkg/frame"{{end}}
)

{{.Doc}}
{{if .FromTop}}
func {{.FuncName}}(req *{{.ReqType}}, needRsp bool, pids ...int32) (rsppara *{{.RspType}}, erro *errorMsg.ErrRsp) {
	return natsHandle.RequestFromTop[{{.ReqType}}, {{.RspType}}](
		"{{.ServerName}}",
		"{{.SFuncName}}",
		req,
		needRsp,
		pids,
	)
}
{{else}}
func {{.FuncName}}(sess *frame.Session,req *{{.ReqType}}, needRsp bool, pids ...int32) (rsppara *{{.RspType}}, erro *errorMsg.ErrRsp) {
	return natsHandle.RequestWithSess[{{.ReqType}}, {{.RspType}}](
		sess,
		"{{.ServerName}}",
		"{{.SFuncName}}",
		req,
		needRsp,
		pids,
	)
}
{{end}}