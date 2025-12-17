package router

import (
	"go.uber.org/zap"

	"testtask/internal/transport/http/handler"
	"testtask/internal/transport/http/middleware"

	"github.com/gin-gonic/gin"
)

type Router struct {
	rout *gin.Engine
	h    *handler.Handler
	log  *zap.Logger
}

func NewRouter(h *handler.Handler, mode string, log *zap.Logger) *Router {
	switch mode {
	case "debug":
		gin.SetMode(gin.DebugMode)
	default:
		gin.SetMode(gin.ReleaseMode)
	}
	router := &Router{
		rout: gin.Default(),
		h:    h,
		log:  log.Named("router"),
	}
	router.setupRouter()

	return router
}

func (r *Router) setupRouter() {
	r.rout.Use(middleware.LoggingMiddleware(r.log))

	gr := r.rout.Group("/")
	r.addApi(gr)

}
func (r *Router) addApi(rg *gin.RouterGroup) {
	api := r.rout.Group("/api/v1")

	api.GET("/wallets/:id", r.h.GetBalance)
	api.POST("/wallet", r.h.Operation)
	api.POST("/wallets", r.h.CreateWallet)
}

func (r *Router) GetEngine() *gin.Engine {
	return r.rout
}

func (r *Router) Start(addr string) error {
	return r.rout.Run(addr)
}
