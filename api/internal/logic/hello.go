package logic

import (
	"context"
	"fmt"
	"go-zero-gin-template/api/internal/svc"
	"go-zero-gin-template/api/internal/types"
)

type HelloLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHelloLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HelloLogic {
	return &HelloLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (h *HelloLogic) Hello(req *types.InfoReq) (string, error) {
	return fmt.Sprintf("hello, %s", req.Name), nil
}
