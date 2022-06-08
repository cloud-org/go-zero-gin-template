package handler

import (
	"go-zero-gin-template/api/internal/svc"
	"go-zero-gin-template/core/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func addApi(r *gin.Engine, serverCtx *svc.ServiceContext) {

	v1 := r.Group("/v1")
	{
		h := Handler{svcCtx: serverCtx}
		helloHandler := HelloHandler{&h}
		helloHandler.InitRouter(v1)
	}

}

// 增加 cors 中间件
// todo: nginx 部署的时候貌似有点问题
func addMiddlewareCors(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"POST", "OPTIONS", "GET", "PUT", "DELETE", "PATCH", "PATCH"},
		AllowHeaders: []string{
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"X-CSRF-Token",
			"Authorization",
			"accept",
			"origin",
			"Cache-Control",
			"X-Requested-With",
		},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
}

//CreateEngine gin
func CreateEngine(serverCtx *svc.ServiceContext) (*gin.Engine, error) {
	if !serverCtx.Config.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(middleware.LogxLogger())
	r.Use(gin.Recovery())
	addMiddlewareCors(r)
	addApi(r, serverCtx)

	return r, nil
}
