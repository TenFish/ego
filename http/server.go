package http

import (
	"errors"
	"github.com/ebar-go/ego/app"
	"github.com/ebar-go/ego/http/handler"
	"github.com/gin-gonic/gin"
	"net"
	"strconv"
	"sync"
)

// Server Web服务管理器
type Server struct {
	// 并发锁,可导出结构体采用私有变量，而不采用内嵌的方式
	mu sync.Mutex

	// gin的路由
	Router *gin.Engine

	// not found handler
	NotFoundHandler func(ctx *gin.Context)
}

// NewServer 实例化server
func NewServer() *Server {
	router := gin.Default()

	return &Server{
		Router:          router,
		NotFoundHandler: handler.NotFoundHandler,
	}
}

// Start run http server
// args must be less than one,
// eg: Start()
//	   Start(8080)
func (server *Server) Start(args ...int) error {
	// use lock
	server.mu.Lock()

	server.beforeStart()

	port := app.Config().ServicePort
	if len(args) == 1 {
		port = args[0]
		app.Config().ServicePort = port
	} else if len(args) > 1 {
		return errors.New("args must be less than one")
	}

	// 404
	server.Router.NoRoute(server.NotFoundHandler)
	server.Router.NoMethod(server.NotFoundHandler)

	completeHost := net.JoinHostPort("", strconv.Itoa(port))

	return server.Router.Run(completeHost)
}

// beforeStart
func (server *Server) beforeStart() {
	// before start
	eventDispatcher := app.EventDispatcher()
	eventDispatcher.Trigger(app.ConfigInitEvent, nil)
	eventDispatcher.Trigger(app.LogManagerInitEvent, nil)

	// mysql auto connect
	if app.Config().Mysql().AutoConnect {
		eventDispatcher.Trigger(app.MySqlConnectEvent, nil)
	}

	// redis auto connect
	if app.Config().Redis().AutoConnect {
		eventDispatcher.Trigger(app.RedisConnectEvent, nil)
	}
}

