package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const timeout = 15 * time.Second

// NewClient established connection to a mongoDb instance using provided URI and auth credentials.
func NewClient(connectionString string) (*mongo.Client, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), timeout)
	defer cancelFunc()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func IsDuplicate(err error) bool {
	var e mongo.WriteException
	if errors.As(err, &e) {
		for _, we := range e.WriteErrors {
			if we.Code == 11000 {
				return true
			}
		}
	}

	return false
}

type Counter struct {
	ID  string `bson:"_id"`
	Seq int    `bson:"seq"`
}

func GetNextSequence(ctx context.Context, counterCollection *mongo.Collection, collectionName string) (int, error) {
	var counter Counter
	filter := bson.M{"_id": collectionName}
	update := bson.M{"$inc": bson.M{"seq": 1}}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	err := counterCollection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&counter)
	if err != nil {
		return 0, err
	}

	return counter.Seq, nil
}

type IndexType int32

const (
	ASC  IndexType = 1
	DESC IndexType = -1

	ConnectionTimeoutInSecond = 10
)

// MustCreateIndex will panic if creating an index on given collection fail
func MustCreateIndex(index mongo.IndexModel, c *mongo.Collection) {
	opts := options.CreateIndexes().SetMaxTime(10 * time.Hour)
	ctx, cancelFunc := context.WithTimeout(context.Background(), ConnectionTimeoutInSecond*time.Hour)
	defer cancelFunc()
	if _, err := c.Indexes().CreateOne(ctx, index, opts); err != nil {
		panic(fmt.Sprintf("error while applying index to collection[%s], error[%s]", c.Name(), err.Error()))
	}
	log.Printf("index[%v] created on collection[%s]", index.Keys, c.Name())
}
