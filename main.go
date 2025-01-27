package main

import (
	"Fire/cache"
	"Fire/config"
	"Fire/pkg/snowflake"
	"Fire/pkg/util"
	log2 "Fire/pkg/util/log"
	"Fire/router"
	"context"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// 1. 加载配置
	if err := config.Init(); err != nil {
		fmt.Printf("config.Init failed, err: %v\n", err)
		return
	}
	if err := snowflake.Init(config.Config.System.StartTime, config.Config.System.MachineID); err != nil {
		fmt.Printf("init snowflake failed, err: %v\n", err)
		return
	}
	cache.InitRedis()
	log2.InitLog()
	//err := util.InitMinio(config.Config.Minio.Endpoint, config.Config.Minio.AccessKey, config.Config.Minio.SecretKey)
	//if err != nil {
	//	fmt.Printf("init minio failed, err: %v\n", err)
	//	return
	//}

	if err := util.InitMinio(config.Config.Minio.Endpoint, config.Config.Minio.AccessKey, config.Config.Minio.SecretKey); err != nil {
		fmt.Printf("init minio failed, err: %v\n", err)
		return
	}

	r := router.NewRouter() // gin路由

	_ = r.Run(config.Config.System.HttpPort)

	router.StartVideoStreamServer() //socket路由

	// 6. 启动服务（优雅关机）
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("port")),
		Handler: r,
	}

	go func() {
		// 开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	zap.L().Info("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown: ", zap.Error(err))
	}

	zap.L().Info("Server exiting")

}
