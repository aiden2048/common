package {{.PkgName}}

import (
	"context"
	"github.com/aiden2048/pkg/qgframe"
	"github.com/aiden2048/pkg/qgframe/logs"
	{{if .HasRequest}}jsoniter "github.com/json-iterator/go"{{end}}
	{{if .HadErr}}"github.com/aiden2048/pkg/public/errorMsg"{{end}}
	{{.ImportPackages}}
)
// {{.Doc}}
func {{.HandlerName}}(r *qgframe.{{.Msg}}) int32 {

	if r == nil {
		logs.LogError("req is nil")
		return qgframe.ESMR_FAILED
	}
	
	{{if .IsEvent}}
		// 解析请求
		{{if .HasRequest}}var req types.{{.RequestType}}
			_err := jsoniter.Unmarshal(r.MsgBody, &req)
			if _err != nil {
				return qgframe.ESMR_FAILED
			}
		{{end}}
		ctx := context.Background()
		l := {{.LogicName}}.New{{.LogicType}}Obj(ctx,r)
		{{if .HasResp}}_, {{end}}err := l.{{.Call}}({{if .HasRequest}}&req{{end}})
		if err != nil {
			return qgframe.ESMR_FAILED
		} 
	{{else}}
		// 解析请求
		{{if .HasRequest}}var req types.{{.RequestType}}
		_err := jsoniter.Unmarshal(r.GetParam(), &req)
		if _err != nil {
			err := errorMsg.Param_Err
			r.SendResponse(err.Ret, "", err.Params)
			return qgframe.ESMR_FAILED
		}
		{{end}}
		ctx := context.Background()
		l := {{.LogicName}}.New{{.LogicType}}Obj(ctx,r)
		{{if .HasResp}}resp, {{end}}err := l.{{.Call}}({{if .HasRequest}}&req{{end}})
		if err != nil {
			r.SendResponse(err.Ret, "", err.Params)
			return qgframe.ESMR_FAILED
		} 

		{{if .HasResp}}
		if qgframe.IsTestUid(r.GetUid()) {
			r.SendResponseEcrypt(0, "", resp)
		} else {
			r.SendResponse(0, "", resp)
		}
		{{end}}
	{{end}}
	return qgframe.ESMR_SUCCEED
}
