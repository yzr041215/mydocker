package deamon

import (
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
)

type Router struct {
	l net.Listener
	r *gin.Engine
}

func NewRouter(l net.Listener) *Router {

	return &Router{
		l: l,
		r: gin.Default(),
	}
}
func (r *Router) Run() error {
	return http.Serve(r.l, r.r.Handler())
}
