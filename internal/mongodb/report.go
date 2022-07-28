package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ITA-Dnipro/Dp-230-Result-Collector/model"
	"github.com/pkg/errors"
)

const (
	reportsDB         = "reports"
	reportsCollection = "reports"
)

var (
	ErrObjectIDTypeConversion = errors.New("object id type conversion")
)

type reportMongoRepo struct {
	mongoDB *mongo.Client
}

func NewReportMongoRepo(mongoDB *mongo.Client) *reportMongoRepo {
	return &reportMongoRepo{mongoDB: mongoDB}
}

func (r *reportMongoRepo) Create(ctx context.Context, report *model.Report) (*model.Report, error) {

	collection := r.mongoDB.Database(reportsDB).Collection(reportsCollection)

	result, err := collection.InsertOne(ctx, report, &options.InsertOneOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "InsertOne")
	}

	objectID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errors.Wrap(ErrObjectIDTypeConversion, "report.InsertedID")
	}

	report.ID = objectID

	return report, nil
}

func (r *reportMongoRepo) PushResult(ctx context.Context, id string, tr model.TestResult) (*model.Report, error) {
	collection := r.mongoDB.Database(reportsDB).Collection(reportsCollection)

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.Wrap(ErrObjectIDTypeConversion, "report.UpdatedID")
	}

	ops := options.FindOneAndUpdate()
	ops.SetReturnDocument(options.After)
	ops.SetUpsert(false)

	upd := bson.M{"$push": bson.M{"testResults": tr}, "$inc": bson.M{"finishTestCount": 1}}

	var report model.Report
	if err := collection.FindOneAndUpdate(ctx, bson.M{"_id": objectID}, upd, ops).Decode(&report); err != nil {
		return nil, errors.Wrap(err, "Decode")
	}

	return &report, nil
}
