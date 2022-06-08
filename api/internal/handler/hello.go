package handler

import (
	"go-zero-gin-template/api/internal/logic"
	"go-zero-gin-template/api/internal/types"

	"github.com/gin-gonic/gin"
)

type HelloHandler struct {
	*Handler
}

func (h *HelloHandler) InitRouter(router *gin.RouterGroup) {
	fn := logic.NewHelloLogic

	router.GET("/hello", h.ReflectHandler(
		fn,
		&types.InfoReq{},
		"Hello",
	))
}
