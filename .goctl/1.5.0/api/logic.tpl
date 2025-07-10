package {{.pkgName}}

import (
	{{if not .hasEvent}}"github.com/aiden2048/pkg/qgframe"{{end}}
	"github.com/aiden2048/pkg/public/errorMsg"
	{{if .hasEvent}}"github.com/aiden2048/common/public/events"
	"context"{{end}}
	{{.imports}}
)

type {{.logic}} struct {

	ctx    context.Context
	{{if not .hasEvent}}r *qgframe.{{.msg}}{{end}}
}

func New{{.logic}}Obj({{if not .hasEvent}}ctx context.Context, r *qgframe.{{.msg}}{{end}}) *{{.logic}} {
	{{if not .hasEvent}}return &{{.logic}}{
		ctx:    ctx,
		r: r,
	}{{else}}return &{{.logic}}{
		ctx: context.TODO(),
	}{{end}}
}
// {{.doc}}
func (l *{{.logic}}) {{.function}}({{.request}}) {{.responseType}} {
	{{.resp}}
	// todo: add your logic here and delete this line

	{{.returnString}}
}
