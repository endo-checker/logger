package store

import (
	"context"

	st "github.com/endo-checker/common/store"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	loggerv1 "github.com/endo-checker/logger/internal/gen/logger/v1"
)

type Storer interface {
	st.Storer[*loggerv1.Log]
	Fetch(ctx context.Context, qr *loggerv1.QueryRequest) ([]*loggerv1.Log, int64, error)
}

type LoggerStore struct {
	*st.Store[*loggerv1.Log]
}

func (s LoggerStore) Fetch(ctx context.Context, qr *loggerv1.QueryRequest) ([]*loggerv1.Log, int64, error) {
	filter := bson.M{}

	if qr.LogId != "" {
		filter = bson.M{"$and": bson.A{filter,
			bson.M{"$or": bson.A{
				bson.M{"logid": primitive.Regex{Pattern: qr.LogId, Options: "i"}},
			}},
		}}
	}

	f := st.WithFilter(filter)

	fo := st.WithFindOptions(options.FindOptions{
		Skip:  &qr.Offset,
		Limit: &qr.Limit,
		Sort:  bson.M{"risk": -1},
	})

	return s.List(ctx, f, fo)
}
