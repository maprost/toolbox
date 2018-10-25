package net

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Server wraps the gin.Engine
type Server struct {
	engine *gin.Engine
	cfg    *Config

	Extension interface{}
}

// NewServer creates a new router
func NewServer(config *Config) *Server {
	router := gin.New()
	router.Use(gin.Recovery())

	return &Server{
		engine:    router,
		cfg:       config,
		Extension: nil,
	}
}

var server *Server
var serverCreateMutex = &sync.Mutex{}

// NewSingleServer creates a singleton 'server'
func NewSingletonServer(creatorFunc func() *Server) *Server {
	serverCreateMutex.Lock()
	defer serverCreateMutex.Unlock()

	if server != nil {
		return server
	}

	server = creatorFunc()
	return server
}

type HandlerFunc func(*Connection)
type WebSocketFunc func(msg []byte, context interface{}, ws WebSocketWriteChannel) bool
type CheckFunc func(*Connection) bool

// Post request method
func (s *Server) Post(relativePath string, handlers HandlerFunc) gin.IRoutes {
	return s.engine.POST(relativePath, s.handlerFunc(handlers))
}

// Get request method
func (s *Server) Get(relativePath string, handlers HandlerFunc) gin.IRoutes {
	return s.engine.GET(relativePath, s.handlerFunc(handlers))
}

// Delete request method
func (s *Server) Delete(relativePath string, handlers HandlerFunc) gin.IRoutes {
	return s.engine.DELETE(relativePath, s.handlerFunc(handlers))
}

// Put request method
func (s *Server) Put(relativePath string, handlers HandlerFunc) gin.IRoutes {
	return s.engine.PUT(relativePath, s.handlerFunc(handlers))
}

// StaticFiles
func (s *Server) StaticFiles(path string, fs http.FileSystem) gin.IRoutes {
	return s.engine.StaticFS(path, fs)
}

// Run the server
func (s *Server) Run() error {
	return s.engine.Run(s.cfg.Host + ":" + s.cfg.Port)
}

// ServeHTTP is to conform to the http.Handler interface
func (s *Server) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	s.engine.ServeHTTP(rw, req)
}

func (s *Server) handlerFunc(requestFunc HandlerFunc) gin.HandlerFunc {
	return func(gin *gin.Context) {
		con := &Connection{
			server: server,
			gin:    gin,
			start:  time.Now(),
			cfg:    s.cfg,
		}
		con.loadDefaultCookies()
		requestFunc(con)
	}
}
