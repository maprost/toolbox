package net

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Router wraps the gin.Engine
type Router struct {
	engine *gin.Engine
	Config *Config
}

// NewRouter creates a new router
func NewRouter() *Router {
	router := gin.New()
	router.Use(gin.Recovery())

	return &Router{
		engine: router,
		Config: &Config{},
	}
}

type HandlerFunc func(*Connection)
type CheckFunc func(*Connection) bool

// Post request method
func (r *Router) Post(relativePath string, handlers HandlerFunc) gin.IRoutes {
	return r.engine.POST(relativePath, netRequest(handlers, r.Config))
}

// Get request method
func (r *Router) Get(relativePath string, handlers HandlerFunc) gin.IRoutes {
	return r.engine.GET(relativePath, netRequest(handlers, r.Config))
}

// Delete request method
func (r *Router) Delete(relativePath string, handlers HandlerFunc) gin.IRoutes {
	return r.engine.DELETE(relativePath, netRequest(handlers, r.Config))
}

// Put request method
func (r *Router) Put(relativePath string, handlers HandlerFunc) gin.IRoutes {
	return r.engine.PUT(relativePath, netRequest(handlers, r.Config))
}

// StaticFiles
func (r *Router) StaticFiles(path string, fs http.FileSystem) gin.IRoutes {
	return r.engine.StaticFS(path, fs)
}

// Run the server
func (r *Router) Run(host string, port string) error {
	return r.engine.Run(host + ":" + port)
}

// ServeHTTP is to conform to the http.Handler interface
func (r *Router) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	r.engine.ServeHTTP(rw, req)
}

func netRequest(requestFunc HandlerFunc, cfg *Config) gin.HandlerFunc {
	return func(gin *gin.Context) {
		net := &Connection{
			gin:   gin,
			start: time.Now(),
			cfg:   cfg,
		}

		net.loadDefaultCookies()
		requestFunc(net)
	}
}
