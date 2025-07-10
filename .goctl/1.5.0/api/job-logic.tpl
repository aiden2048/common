package {{.pkgName}}

import (
	"github.com/aiden2048/pkg/frame"
	"github.com/aiden2048/pkg/public/errorMsg"
	"github.com/aiden2048/common/serviceApi/jobSchedApi"
	{{if .hasEvent}}"github.com/aiden2048/common/public/events"{{end}}
	{{.imports}}
)

type {{.logic}} struct {
	ctx    context.Context
	r *frame.{{.msg}}
}

func New{{.logic}}Obj(ctx context.Context, r *frame.{{.msg}}) *{{.logic}} {
	return &{{.logic}}{
		ctx:    ctx,
		r: r,
	}
}

// {{.doc}}
func (l *{{.logic}}) {{.function}}(req *jobSchedApi.ExecJobReq) (resp *jobSchedApi.ExecJobResp, err *errorMsg.ErrRsp) {
    // 耗时作作业请更换 jobSchedApi.AsyncExecJob 实现！！因为调度运行作业总线程数为40个，避免阻塞调度器可用运行线程
	return jobSchedApi.ExecJob(req, {{.jobArg}}, func({{.request}}) (extra string, err *errorMsg.ErrRsp) {
    		return "", nil
    	})
}
