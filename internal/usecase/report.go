package usecase

import (
	"context"

	"github.com/ITA-Dnipro/Dp-230-Result-Collector/model"
)

type ReportRepo interface {
	Create(ctx context.Context, report *model.Report) (*model.Report, error)
	PushResult(ctx context.Context, id string, tr model.TestResult) (*model.Report, error)
}

type Validate interface {
	StructCtx(context.Context, interface{}) error
}

type Producer interface {
	Send(r *model.Report) error
}

type reportUsecase struct {
	repo     ReportRepo
	validate Validate
	producer Producer
}

func NewReportUsecase(repo ReportRepo, validate Validate, producer Producer) *reportUsecase {
	return &reportUsecase{repo: repo, validate: validate, producer: producer}
}

func (u *reportUsecase) Create(ctx context.Context, report *model.Report) (*model.Report, error) {
	if err := u.validate.StructCtx(ctx, report); err != nil {
		return nil, err
	}

	return u.repo.Create(ctx, report)

}
func (u *reportUsecase) PushResult(ctx context.Context, id string, tr model.TestResult) (*model.Report, error) {
	report, err := u.repo.PushResult(ctx, id, tr)
	if err != nil {
		return nil, err
	}
	if report.FinishTestCount == report.TotalTestCount {
		u.producer.Send(report)
	}
	return report, nil
}
