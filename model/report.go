package model

import (
	pb "github.com/ITA-Dnipro/Dp-230-Result-Collector/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Report struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	URL             string             `json:"url" bson:"url,omitempty" validate:"required,url"`
	Email           string             `json:"email" bson:"email,omitempty" validate:"required,email"`
	TestResults     []TestResult       `json:"testResults" bson:"testResults,omitempty"`
	TotalTestCount  int64              `json:"total_test_count" bson:"totalTestCount,omitempty" validate:"required"`
	FinishTestCount int64              `json:"finish_test_count" bson:"finishTestCount"`
}

func ReportFromProto(report *pb.Report) (*Report, error) {
	reportID, err := primitive.ObjectIDFromHex(report.GetID())
	if err != nil {
		return nil, err
	}

	rep := &Report{
		ID:    reportID,
		URL:   report.URL,
		Email: report.Email,
	}
	for _, tr := range report.GetTestResults() {
		rep.TestResults = append(rep.TestResults, TestResultFromProto(tr))
	}
	return rep, nil
}

func (rp *Report) ToProto() *pb.Report {
	r := &pb.Report{
		ID:              rp.ID.Hex(),
		URL:             rp.URL,
		Email:           rp.Email,
		TotalTestCount:  rp.TotalTestCount,
		FinishTestCount: rp.FinishTestCount,
	}

	for _, tr := range rp.TestResults {
		r.TestResults = append(r.TestResults, tr.ToProto())
	}
	return r
}
