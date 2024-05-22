package main

import (
	"chat_agent/config"
	"chat_agent/logger"
	"chat_agent/profile"
	"compress/gzip"
	"fmt"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"time"
)

type Addr struct {
	IP   string
	Port int
}

type ServiceConfig struct {
	http       Addr
	httpServer *http.Server

	wg            sync.WaitGroup
	fsnotifyMonit *profile.FsnotifyMonit
}

// 相关服务初始化
func pkgInit() {
	config.LoadConfig()
	//logger.Info("pkg init succ")
}

func main() {
	var cstZone = time.FixedZone("CST", 8*3600) // 东八
	time.Local = cstZone
	service := ServiceConfig{
		http: Addr{IP: config.Server.Http.Ip, Port: config.Server.Http.Port},
	}
	//开始监听
	service.pkgWatcher(pkgInit)
	service.run()
}

// 配置动态监听加载
func (server *ServiceConfig) pkgWatcher(fn func()) {
	watcher, err := profile.NewFsnotify()
	server.fsnotifyMonit = watcher
	//defer watcher.Close()
	if err != nil {
		logger.Fatalf("pkgWatcher error %v", err)
	}
	configPath := config.GetConfigPath()
	if err := watcher.AddWatcher(configPath); err != nil {
		logger.Fatalf("pkgWatcher add path error %v", err)
	}
	logger.Info("pkgWatcher configPath:", configPath)
	go watcher.Run(fn)
}

func (server *ServiceConfig) run() {
	server.serveHttp()
	server.serveGrpc()
	//server.serveAsynq()
	server.gracefulTerminate()
}

func (server *ServiceConfig) serveHttp() {
	var err error
	g := gin.New()
	g.Use(logging.RegisterLog())
	//g.Use(gin.RecoveryWithWriter(logging.WriterLevel(logger.ErrorLevel)))
	g.Use(logging.Recovery(logging.RecoveryHandler))
	g.Use(author.CheckSign())
	g.Use(author.JWTAuth())
	g.Use(gzip.Gzip(gzip.DefaultCompression))
	//g.Use(gin2.Tracer())

	g.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"code": "404",
			"msg":  "ErrPageNotFoundRequest",
			"body": map[string]interface{}{},
		})
	})
	g.NoMethod(func(c *gin.Context) {
		c.JSON(405, gin.H{
			"code": "405",
			"msg":  "ErrMethodNotAllowRequest",
			"body": map[string]interface{}{},
		})
	})
	routers.RegisterMainRouters(g)

	//new server
	addr := fmt.Sprintf("%s:%d", server.http.IP, server.http.Port)
	server.httpServer = &http.Server{
		Addr:    addr,
		Handler: g,
	}
	go func() {
		server.wg.Add(2)
		if err = server.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(fmt.Sprintf("[HTTP 服务结束] err:%s", err.Error()))
		} else {
			logger.Info(fmt.Sprintf("[HTTP 服务结束]"))
		}
		server.wg.Done()
	}()
}
