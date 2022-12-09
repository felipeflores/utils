package persistence

import "context"

type MongoProvider interface {
	InsertOne(ctx context.Context, collection string, t interface{})
	Aggregate(ctx context.Context, collection string, t interface{}, results interface{}) error
}
