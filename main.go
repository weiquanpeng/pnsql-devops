// main.go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/xiaohongshu/PnSql/server/global"
	"github.com/xiaohongshu/PnSql/server/initialize"
)

func main() {
	initialize.InitConfig()

	// 获取采集间隔（示例：从配置中读取，这里硬编码为10秒）
	collectInterval := 10 * time.Second
	collector := initialize.InitSimpleCollector(collectInterval)
	collector.Start()
	log.Println("采集器启动，将每", collectInterval, "打印一次 hahaah")

	router := initialize.InitRouter()
	server := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%s", global.P_cfg.System.Port),
		Handler: router,
	}

	// 启动HTTP服务器
	go func() {
		log.Printf("HTTP服务器启动，监听端口: %s", global.P_cfg.System.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP服务器启动失败: %v", err)
		}
	}()

	// 等待终止信号
	waitForShutdown(server, collector)
}

// 等待关闭并优雅停止服务
func waitForShutdown(server *http.Server, collector *initialize.SimpleCollector) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("\n接收到关闭信号，开始停止服务...")
	// 停止采集器
	collector.Stop()
	// 优雅关闭HTTP服务器
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("服务器关闭失败: %v", err)
	}
	log.Println("所有服务已安全停止")
	os.Exit(0)
}
