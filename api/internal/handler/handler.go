package handler

import (
	"bytes"
	"errors"
	"fmt"
	"go-zero-gin-template/api/internal/svc"
	"go-zero-gin-template/core/common"
	"net/http"
	"reflect"

	"github.com/go-playground/validator/v10"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svcCtx *svc.ServiceContext
}

//ReflectHandler 反射通用 handler
//fn -> NewTargetLogic 实例化 logic 的函数签名
//req -> 请求体 需要传指针类型 .Elem 需要 Kind 为 Ptr
//callFn -> 调用的函数名称
func (h *Handler) ReflectHandler(fn, req interface{}, callFn string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// bind req
		reqType := reflect.TypeOf(req).Elem()
		request := reflect.New(reqType)
		if err := c.ShouldBind(request.Interface()); err != nil {
			// 更显式的字段错误提醒
			var ve validator.ValidationErrors
			if errors.As(err, &ve) {
				errMsg := h.handleBindErr(ve)
				c.JSON(http.StatusOK, common.RespWithError(errMsg))
				return
			}
			c.JSON(http.StatusOK, common.RespWithErrorf("反序列化失败: %v", err))
			return
		}
		logx.Infof("req: %+v", request.Interface())
		// 关键点: 实例化 logic
		logicParams := []reflect.Value{reflect.ValueOf(c), reflect.ValueOf(h.svcCtx)}
		rets := reflect.ValueOf(fn).Call(logicParams)
		if len(rets) != 1 {
			c.JSON(http.StatusOK, fmt.Sprintf("实例化 logic 失败"))
			return
		}
		//  Logic.Call func by method name
		reqParams := []reflect.Value{reflect.ValueOf(request.Interface())}
		method := rets[0].MethodByName(callFn)
		if !method.IsValid() { // if not method
			c.JSON(http.StatusOK, common.RespWithErrorf("%+v 没有实现对应的 method: %v", rets[0].String(), callFn))
			return
		}
		res := method.Call(reqParams)
		if len(res) != 2 {
			c.JSON(http.StatusOK, common.RespWithErrorf("调用 %s func 失败", callFn))
			return
		}

		if !res[1].IsNil() { // 第二个返回参数是 error
			logx.Errorf(
				"fn: %+v, req: %+v, callFn: %s, err: %+v",
				rets[0].String(),
				request.Interface(),
				callFn,
				res[1].Interface(),
			)
			c.JSON(http.StatusOK, common.RespWithErrorf(
				"%+v", res[1].Interface()),
			)
			return
		}

		c.JSON(http.StatusOK, common.RespWithData(res[0].Interface())) // 第一个返回参数是 data
		return
	}
}

func (h *Handler) handleBindErr(validatorErrors validator.ValidationErrors) string {
	var errMsg bytes.Buffer
	for i := 0; i < len(validatorErrors); i++ {
		// 目前只需要处理 required 即可 然而这里还是只能列出对应结构体的字段而已 暂时列不出来对应的 tag 的字段
		// 例如 Id form:"id" 这里只会是 Id 而不是 id 有些可能还是会对应不上
		switch validatorErrors[i].Tag() {
		case "required":
			errMsg.WriteString(fmt.Sprintf("字段: %v 是必填的", validatorErrors[i].Field()))
		default:
			errMsg.WriteString(fmt.Sprintf(
				"字段: %v 要求符合 %v %v 规则",
				validatorErrors[i].Field(),
				validatorErrors[i].Tag(),
				validatorErrors[i].Param(), // 假如 gt=3 这里会是 3
			))
		}
	}

	return errMsg.String()
}
