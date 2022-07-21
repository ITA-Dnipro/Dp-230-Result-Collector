package model

import pb "github.com/ITA-Dnipro/Dp-230-Result-Collector/proto"

// PoC is PoC struct for Result
type PoC struct {
	Type       string `json:"type" bson:"type,omitempty"`
	InjectType string `json:"inject_type" bson:"injectType,omitempty"`
	PoCType    string `json:"poc_type" bson:"pocType,omitempty"`
	Method     string `json:"method" bson:"method,omitempty"`
	Data       string `json:"data" bson:"data,omitempty"`
	Param      string `json:"param" bson:"param,omitempty"`
	Payload    string `json:"payload" bson:"Payload,omitempty"`
	Evidence   string `json:"evidence" bson:"Evidence,omitempty"`
	CWE        string `json:"cwe" bson:"cwe,omitempty"`
	Severity   string `json:"severity" bson:"Severity,omitempty"`
}

func (p *PoC) ToProto() *pb.PoC {
	return &pb.PoC{
		Type:       p.Type,
		InjectType: p.InjectType,
		PoCType:    p.PoCType,
		Method:     p.Method,
		Data:       p.Data,
		Param:      p.Param,
		Payload:    p.Payload,
		Evidence:   p.Evidence,
		SWE:        p.CWE,
		Severity:   p.Severity,
	}
}

func PoCFromProto(p *pb.PoC) PoC {
	return PoC{
		Type:       p.GetType(),
		InjectType: p.GetInjectType(),
		PoCType:    p.GetPoCType(),
		Method:     p.GetMethod(),
		Data:       p.GetData(),
		Param:      p.GetParam(),
		Payload:    p.GetPayload(),
		Evidence:   p.GetEvidence(),
		CWE:        p.GetSWE(),
		Severity:   p.GetSeverity(),
	}
}
