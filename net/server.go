package net

import (
	"net/http"
	"time"

	"sync"

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
type CheckFunc func(*Connection) bool

// Post request method
func (s *Server) Post(relativePath string, handlers HandlerFunc) gin.IRoutes {
	return s.engine.POST(relativePath, netRequest(handlers, s))
}

// Get request method
func (s *Server) Get(relativePath string, handlers HandlerFunc) gin.IRoutes {
	return s.engine.GET(relativePath, netRequest(handlers, s))
}

// Delete request method
func (s *Server) Delete(relativePath string, handlers HandlerFunc) gin.IRoutes {
	return s.engine.DELETE(relativePath, netRequest(handlers, s))
}

// Put request method
func (s *Server) Put(relativePath string, handlers HandlerFunc) gin.IRoutes {
	return s.engine.PUT(relativePath, netRequest(handlers, s))
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

func netRequest(requestFunc HandlerFunc, s *Server) gin.HandlerFunc {
	return func(gin *gin.Context) {
		con := &Connection{
			gin:   gin,
			start: time.Now(),
			cfg:   s.cfg,
		}
		con.loadDefaultCookies()

		if s.cfg.InitConnection != nil {
			err := s.cfg.InitConnection(s, con)
			if err != nil {
				panic(err)
			}
		}

		requestFunc(con)
	}
}
