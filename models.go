package main

import "time"

// mongodb request and response models
type MongoRequest struct {
	MinCount  int    `json:"minCount"`
	MaxCount  int    `json:"maxCount"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

type MongoResponse struct {
	Code    int                    `json:"code"`
	Msg     string                 `json:"msg"`
	Records *[]MongoResponseRecord `json:"records"`
}

type MongoResponseRecord struct {
	Key           string    `json:"key" bson:"key"`
	CreatedAt     string    `json:"createdAt"`
	TotalCount    int       `json:"totalCount" bson:"totalCount"`
	CreatedAtTime time.Time `json:"-" bson:"createdAt"`
}

//in-memory db models
type InMemory struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// other
type CustomError struct {
	Error string `json:"error"`
}
