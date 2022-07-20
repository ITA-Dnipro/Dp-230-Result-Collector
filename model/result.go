package model

import (
	"time"

	pb "github.com/ITA-Dnipro/Dp-230-Result-Collector/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Result struct {
	URL       string    `json:"url" bson:"url,omitempty"`
	PoCs      []PoC     `json:"pocs" bson:"pocs,omitempty"`
	StartTime time.Time `json:"start_time" bson:"startTime,omitempty" `
	EndTime   time.Time `json:"end_time" bson:"endTime,omitempty"`
}

type TestResult struct {
	Type    string   `json:"type" bson:"type,omitempty"`
	Results []Result `json:"results" bson:"results,omitempty"`
}

func ResultFromProto(res *pb.Result) Result {
	r := Result{
		URL:       res.GetURL(),
		StartTime: res.GetStartTime().AsTime(),
		EndTime:   res.GetEndTime().AsTime(),
	}
	for _, poc := range res.GetPoCs() {
		r.PoCs = append(r.PoCs, PoCFromProto(poc))
	}
	return r
}

func TestResultFromProto(tr *pb.TestResult) TestResult {
	r := TestResult{
		Type: tr.Type,
	}
	for _, res := range tr.GetResults() {
		r.Results = append(r.Results, ResultFromProto(res))
	}
	return r
}

func (res *Result) ToProto() *pb.Result {
	r := &pb.Result{
		URL:       res.URL,
		StartTime: timestamppb.New(res.StartTime),
		EndTime:   timestamppb.New(res.EndTime),
	}

	for _, p := range res.PoCs {
		r.PoCs = append(r.PoCs, p.ToProto())
	}
	return r
}
func (tr *TestResult) ToProto() *pb.TestResult {
	res := &pb.TestResult{
		Type: tr.Type,
	}

	for _, r := range tr.Results {
		res.Results = append(res.Results, r.ToProto())
	}
	return res
}
