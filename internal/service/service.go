package service

import (
	"context"
	"fmt"

	"github.com/ITA-Dnipro/Dp-230-Result-Collector/model"
	pb "github.com/ITA-Dnipro/Dp-230-Result-Collector/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type reportRepo interface {
	Create(ctx context.Context, report *model.Report) (*model.Report, error)
	PushResult(ctx context.Context, id string, result model.Result) (*model.Report, error)
}

type reportService struct {
	repo reportRepo
}

func NewReportService(repo reportRepo) *reportService {
	return &reportService{repo: repo}
}

func (r *reportService) Create(ctx context.Context, req *pb.CreateReq) (*pb.CreateRes, error) {
	report := &model.Report{TotalTestCount: req.GetTotalTestCount()}
	created, err := r.repo.Create(ctx, report)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("%s: %v", err.Error(), err))
	}
	return &pb.CreateRes{Report: created.ToProto()}, nil
}

func (r *reportService) PushResult(ctx context.Context, req *pb.PushResultReq) (*pb.PushResultRes, error) {
	id := req.GetID()
	result := model.ResultFromProto(req.GetResult())
	created, err := r.repo.PushResult(ctx, id, result)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("%s: %v", err.Error(), err))
	}
	return &pb.PushResultRes{Report: created.ToProto()}, nil
}
