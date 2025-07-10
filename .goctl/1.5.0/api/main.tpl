package main

import (
	"log"
	"github.com/aiden2048/pkg/public/mongodb"
	"github.com/aiden2048/pkg/frame"
	"github.com/aiden2048/pkg/utils/baselib"
	{{if .httpServer}}"fmt"
	"github.com/aiden2048/common/public/commonConst"
	"github.com/valyala/fasthttp"{{end}}
	{{.importPackages}}
)


func main() {
	if err := frame.InitConfig("{{.serviceName}}", &frame.FrameOption{
		
	}); err != nil {
		log.Fatalf("InitConfig error:%s", err.Error())
		return
	}
	// 初始化mongo
	if err := mongodb.StartMgoDb(mongodb.WLevel2); err != nil {
		log.Fatalf("InitMongodb %+v failed: %s", frame.GetMgoCoinfig(), err.Error())
		return
	}
	// TODO 加载配置（可选）
	if err := baselib.InitConfig(config.LoadConfig); err != nil {
		log.Fatalf(" baselib.InitConfig failed: %s", err.Error())
		return
	}
	{{if .httpServer}}
	// 启动 HTTP 服务器
	server := fasthttp.Server{
		Handler: handler.RegisterHttpHandlers(),
	}
	go server.ListenAndServe(fmt.Sprintf(":%d", commonConst.HTTP_PORT_{{.upServiceName}})) 
	{{end}}
	defer frame.Stop()
	frame.Run(handler.RegisterHandlers)
}
