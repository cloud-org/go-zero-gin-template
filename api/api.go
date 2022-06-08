package main

import (
	"context"
	"flag"
	"fmt"
	"go-zero-gin-template/api/internal/config"
	"go-zero-gin-template/api/internal/handler"
	"go-zero-gin-template/api/internal/svc"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
	filePath string // 配置文件路径
	help     bool   // 帮助
)

func usage() {
	fmt.Fprintf(os.Stdout, `./api - go-zero-gin server
	Usage: ./api [-h help] [-c etc/dev.yaml]
	Options:
	`)
	flag.PrintDefaults()
}

func main() {
	flag.StringVar(&filePath, "c", "etc/dev.yaml", "配置文件所在")
	flag.BoolVar(&help, "h", false, "帮助")
	flag.Usage = usage
	flag.Parse()
	if help {
		flag.PrintDefaults()
		return
	}

	// 配置文件读取
	var c config.Config
	conf.MustLoad(filePath, &c)
	//log.Printf("config is %+v\n", c)

	// 设置logx
	logx.MustSetup(c.Log)

	svcCtx := svc.NewServiceContext(c)

	address := net.JoinHostPort(svcCtx.Config.Host, fmt.Sprintf("%d", svcCtx.Config.Port))
	if address == "" {
		logx.Error("err listen addr")
		return
	}

	router, err := handler.CreateEngine(svcCtx)
	if err != nil {
		logx.Errorf("create server err: %v", err)
		return
	}

	server := &http.Server{
		Addr:    address,
		Handler: router,
	}

	go func() {
		logx.Infof("server is start at: %v", address)
		if err = server.ListenAndServe(); err != nil {
			if strings.Contains(err.Error(), "bind: address already in use") {
				logx.Error("端口被占用, %s", err.Error())
				panic(err)
				return
			}
			logx.Errorf("listen server err: %v", err) // 这里不要用 Fatal 不然优雅关停会直接退出
		}
	}()

	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	// DONE: 优雅关停
	for {
		s := <-osSignal
		logx.Infof("[main] 捕获信号 %s", s.String())
		switch s {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			if err = server.Shutdown(ctx); err != nil {
				logx.Errorf("[main] 程序退出异常 err:%s", err.Error())
			} else {
				logx.Info("[main] 程序正常退出")
			}
			logx.Close()
			cancel()
			return
		case syscall.SIGHUP:
			logx.Info("SIGHUP, ignore")
		default:
			logx.Info("other signal")
		}
	}
}
