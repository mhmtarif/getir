package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const uri = "mongodb+srv://challengeUser:WUMglwNBaydH8Yvu@challenge-xzwqd.mongodb.net/getircase-study?retryWrites=true"
const layout = "2006-01-02"

var records *mongo.Collection

// var client *mongo.Client

func initMongoDb() {

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	records = client.Database("getir-case-study").Collection("records")
}

func getUnixTime(dateStr string) (time.Time, error) {
	return time.Parse(layout, dateStr)
}

func getDateStr(date time.Time) string {
	return date.Format(layout)
}

func getRecords(startDate, endDate string, minCount, maxCount int) ([]MongoResponseRecord, error) {
	var results []MongoResponseRecord
	start, err := getUnixTime(startDate)
	if err != nil {
		return results, err
	}

	end, err := getUnixTime(endDate)
	if err != nil {
		return results, err
	}

	stage1 := bson.D{{"$match", bson.M{"createdAt": bson.M{
		"$gte": primitive.NewDateTimeFromTime(start),
		"$lte": primitive.NewDateTimeFromTime(end),
	}}}}

	stage2 := bson.D{{"$addFields", bson.M{"totalCount": bson.M{
		"$sum": "$counts",
	}}}}

	stage3 := bson.D{{"$match", bson.M{"totalCount": bson.M{
		"$gte": minCount,
		"$lte": maxCount,
	}}}}

	stage4 := bson.D{{"$project", bson.D{{"createdAt", 1}, {"key", 1}, {"totalCount", 1}, {"_id", 0}}}}

	cursor, err := records.Aggregate(context.TODO(), mongo.Pipeline{stage1, stage2, stage3, stage4})
	if err != nil {
		return results, err

	}

	if err = cursor.All(context.TODO(), &results); err != nil {
		return results, err
	}
	for i := 0; i < len(results); i++ {
		date := results[i].CreatedAtTime
		results[i].CreatedAt = getDateStr(date)
	}
	return results, nil
}
