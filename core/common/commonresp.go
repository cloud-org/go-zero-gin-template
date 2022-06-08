package common

import (
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
)

// 定义错误码
const (
	ErrCodeTokenAuthErr = 11
	ErrCodeTokenExpired = 12
)

// CommonResp ok 则为 0，err 则为 -1
type CommonResp struct {
	ErrMsg  string      `json:"errMsg"`
	ErrCode int         `json:"errCode"`
	Data    interface{} `json:"data"`
}

//RespNoData 正常返回
func RespNoData() *CommonResp {
	return &CommonResp{}
}

//RespWithData 正常返回
func RespWithData(data interface{}) *CommonResp {
	resp := CommonResp{}
	resp.Data = data
	return &resp
}

//RespWithErrorCode 携带自定义错误码
func RespWithErrorCode(errCode int, err string) *CommonResp {
	resp := CommonResp{ErrCode: errCode}
	resp.ErrMsg = err
	logx.Errorf("错误响应: %v", resp.ErrMsg)
	return &resp
}

//RespWithError 错误返回
func RespWithError(err string) *CommonResp {
	resp := CommonResp{ErrCode: -1}
	resp.ErrMsg = err
	logx.Errorf("错误响应: %v", resp.ErrMsg)
	return &resp
}

func RespWithErrorf(format string, a ...interface{}) *CommonResp {
	resp := CommonResp{ErrCode: -1}
	resp.ErrMsg = fmt.Sprintf(format, a...)
	logx.Errorf("错误响应: %v", resp.ErrMsg)
	return &resp
}
