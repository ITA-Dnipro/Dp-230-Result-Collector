package service

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ITA-Dnipro/Dp-230-Result-Collector/model"
	pb "github.com/ITA-Dnipro/Dp-230-Result-Collector/proto"
)

type reportUsecase interface {
	Create(ctx context.Context, report *model.Report) (*model.Report, error)
	PushResult(ctx context.Context, id string, tr model.TestResult) (*model.Report, error)
	GetReport(ctx context.Context, id string) (*model.Report, error)
}

type reportService struct {
	usecase reportUsecase
}

func NewReportService(uc reportUsecase) *reportService {
	return &reportService{usecase: uc}
}

func (r *reportService) Create(ctx context.Context, req *pb.CreateReq) (*pb.CreateRes, error) {
	report := &model.Report{
		URL:            req.GetURL(),
		Email:          req.GetEmail(),
		TotalTestCount: req.GetTotalTestCount(),
	}

	created, err := r.usecase.Create(ctx, report)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("%s: %v", err.Error(), err))
	}

	return &pb.CreateRes{Report: created.ToProto()}, nil
}

func (r *reportService) PushResult(ctx context.Context, req *pb.PushResultReq) (*pb.PushResultRes, error) {
	id := req.GetID()
	tr := model.TestResultFromProto(req.GetTestResult())

	created, err := r.usecase.PushResult(ctx, id, tr)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("%s: %v", err.Error(), err))
	}

	return &pb.PushResultRes{Report: created.ToProto()}, nil
}

func (r *reportService) GetReport(ctx context.Context, req *pb.GetReportReq) (*pb.GetReportRes, error) {
	id := req.GetID()

	report, err := r.usecase.GetReport(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("%s: %v", err.Error(), err))
	}

	return &pb.GetReportRes{Report: report.ToProto()}, nil
}
