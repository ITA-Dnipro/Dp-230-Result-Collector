package model

import (
	pb "github.com/ITA-Dnipro/Dp-230-Result-Collector/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Report struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Results         []Result           `json:"results" bson:"results,omitempty"`
	TotalTestCount  int64              `json:"total_test_count" bson:"totalTestCount,omitempty"`
	FinishTestCount int64              `json:"finish_test_count" bson:"finishTestCount"`
}

func ReportFromProto(report *pb.Report) (*Report, error) {
	reportID, err := primitive.ObjectIDFromHex(report.GetID())
	if err != nil {
		return nil, err
	}

	rep := &Report{
		ID: reportID,
	}
	for _, res := range report.GetResults() {
		rep.Results = append(rep.Results, ResultFromProto(res))
	}
	return rep, nil
}

func (rp *Report) ToProto() *pb.Report {
	r := &pb.Report{
		ID:              rp.ID.String(),
		TotalTestCount:  rp.TotalTestCount,
		FinishTestCount: rp.FinishTestCount,
	}

	for _, res := range rp.Results {
		r.Results = append(r.Results, res.ToProto())
	}
	return r
}
