package servprof

//import (
//	"github.com/gin-gonic/gin"
//	"github.com/mentalisit/logger"
//	"net/http/pprof"
//)
//
//type Server struct {
//	log *logger.Logger
//}
//
//func NewServer(log *logger.Logger) *Server {
//	s := &Server{log: log}
//	go s.runServer()
//	return s
//}
//func (s *Server) runServer() {
//	gin.SetMode(gin.ReleaseMode)
//	r := gin.New()
//
//	pprofRoutes(r)
//
//	err := r.Run(":80")
//	if err != nil {
//		s.log.ErrorErr(err)
//		return
//	}
//}
//
//// pprofRoutes регистрирует обработчики pprof на маршрутизаторе Gin
//func pprofRoutes(router *gin.Engine) {
//	pprofGroup := router.Group("/debug/pprof")
//	{
//		pprofGroup.GET("/", gin.WrapF(pprof.Index))
//		pprofGroup.GET("/cmdline", gin.WrapF(pprof.Cmdline))
//		pprofGroup.GET("/profile", gin.WrapF(pprof.Profile))
//		pprofGroup.POST("/symbol", gin.WrapF(pprof.Symbol))
//		pprofGroup.GET("/symbol", gin.WrapF(pprof.Symbol))
//		pprofGroup.GET("/trace", gin.WrapF(pprof.Trace))
//		pprofGroup.GET("/allocs", gin.WrapH(pprof.Handler("allocs")))
//		pprofGroup.GET("/block", gin.WrapH(pprof.Handler("block")))
//		pprofGroup.GET("/goroutine", gin.WrapH(pprof.Handler("goroutine")))
//		pprofGroup.GET("/heap", gin.WrapH(pprof.Handler("heap")))
//		pprofGroup.GET("/mutex", gin.WrapH(pprof.Handler("mutex")))
//		pprofGroup.GET("/threadcreate", gin.WrapH(pprof.Handler("threadcreate")))
//	}
//}
