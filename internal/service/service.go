package service

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ITA-Dnipro/Dp-230-Result-Collector/model"
	pb "github.com/ITA-Dnipro/Dp-230-Result-Collector/proto"
)

type reportRepo interface {
	Create(ctx context.Context, report *model.Report) (*model.Report, error)
	PushResult(ctx context.Context, id string, tr model.TestResult) (*model.Report, error)
}

type reportService struct {
	repo     reportRepo
	validate *validator.Validate
}

func NewReportService(repo reportRepo, validate *validator.Validate) *reportService {
	return &reportService{repo: repo, validate: validate}
}

func (r *reportService) Create(ctx context.Context, req *pb.CreateReq) (*pb.CreateRes, error) {
	report := &model.Report{
		URL:            req.GetURL(),
		Email:          req.GetEmail(),
		TotalTestCount: req.GetTotalTestCount(),
	}

	if err := r.validate.StructCtx(ctx, report); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("%s: %v", err.Error(), err))
	}

	created, err := r.repo.Create(ctx, report)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("%s: %v", err.Error(), err))
	}

	return &pb.CreateRes{Report: created.ToProto()}, nil
}

func (r *reportService) PushResult(ctx context.Context, req *pb.PushResultReq) (*pb.PushResultRes, error) {
	id := req.GetID()
	tr := model.TestResultFromProto(req.GetTestResult())
	created, err := r.repo.PushResult(ctx, id, tr)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("%s: %v", err.Error(), err))
	}
	return &pb.PushResultRes{Report: created.ToProto()}, nil
}
