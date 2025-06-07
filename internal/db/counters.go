package db

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var initialCounterIDs = []string{"user_id"}

/*
GetNextSequenceValue 获取并自增指定计数器的值，返回自增后的值。
如果计数器不存在，则创建并从1开始。
*/
func GetNextSequenceValue(counterID string) (int64, error) {
	conn, err := GetMongoConnector()
	if err != nil {
		return 0, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": counterID}
	update := bson.M{"$inc": bson.M{"sequence_value": 1}}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	var counterDoc struct {
		SequenceValue int64 `bson:"sequence_value"`
	}
	err = conn.DB.Collection("counters").FindOneAndUpdate(ctx, filter, update, opts).Decode(&counterDoc)
	if errors.Is(err, mongo.ErrNoDocuments) {
		// 计数器不存在，创建并设置为1
		counterDoc.SequenceValue = 1
		_, err = conn.DB.Collection("counters").UpdateOne(
			ctx,
			bson.M{"_id": counterID},
			bson.M{"$set": bson.M{"sequence_value": 1}},
			options.Update().SetUpsert(true),
		)
		if err != nil {
			return 0, err
		}
		return 1, nil
	} else if err != nil {
		return 0, err
	}
	return counterDoc.SequenceValue, nil
}
