package datacenterhttphandler

import (
	datacenterPb "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/datacenter/datacenterPb"
	datacenterusecase "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/datacenter/datacenterUsecase"
	"context"
)

type (
	datacenterGrpcHandler struct {
		datacenterPb.UnimplementedDataCenterGrpcServiceServer
		datacenterUsecase datacenterusecase.DatacenterUsecaseService
	}
)

func NewdatacenterGrpcHandler(datacenterUsecase datacenterusecase.DatacenterUsecaseService) *datacenterGrpcHandler {
	return &datacenterGrpcHandler{datacenterUsecase: datacenterUsecase}
}

func (g *datacenterGrpcHandler) GetProjectDataCenter(ctx context.Context, req *datacenterPb.GetProjectDataCenterReq) (*datacenterPb.GetProjectDataCenterRes, error) {
	return g.datacenterUsecase.GetAllProjectUsecase(ctx, req.Limit, req.Skip)
}

func (g *datacenterGrpcHandler) GetSingleProjectDataCenter(ctx context.Context, req *datacenterPb.GetSingleProjectDataCenterReq) (*datacenterPb.GetSingleProjectDataCenterRes, error) {
	return g.datacenterUsecase.GetSingleProjectUsecase(ctx, req.ProjecId)
}
