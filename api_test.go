package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func postRequest(queryData interface{}, url string, apitype string) io.ReadCloser {
	body, _ := json.Marshal(queryData)
	req, _ := http.NewRequest("POST", url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	w := httptest.NewRecorder()
	if apitype == "mongoapi" {
		mongoApi(w, req)
	} else {
		memoryApi(w, req)
	}
	res := w.Result()
	defer res.Body.Close()
	return res.Body
}

func getMemoryRequest(key string) (InMemory, error) {
	req2, _ := http.NewRequest("GET", "/in-memory?key="+key, nil)
	w2 := httptest.NewRecorder()
	memoryApi(w2, req2)
	res2 := w2.Result()
	decoder2 := json.NewDecoder(res2.Body)
	var resp2 InMemory
	err := decoder2.Decode(&resp2)
	return resp2, err
}

func TestMongoApi(t *testing.T) {
	initMongoDb()
	queryData := &MongoRequest{StartDate: "2016-01-26", EndDate: "2018-02-02", MinCount: 2700, MaxCount: 3000}
	body := postRequest(queryData, "/mongo", "mongoapi")
	decoder := json.NewDecoder(body)
	var resp MongoResponse
	err := decoder.Decode(&resp)
	if err != nil {
		t.Errorf("expected err to be nil got %v", err)
	}
	if resp.Msg != "Success" {
		t.Errorf("expected success got %v", resp.Msg)
	}
	if resp.Code != 0 {
		t.Errorf("expected 0 got %v", resp.Code)
	}

	if len(*resp.Records) == 0 {
		t.Errorf("expected an nonempty list got empty list")
	}

}

func TestMemoryApi(t *testing.T) {
	initMemoryDb()
	queryData := &InMemory{Key: "key1", Value: "val1"}
	body := postRequest(queryData, "/in-memory", "in-memory")
	decoder := json.NewDecoder(body)
	var resp InMemory
	err := decoder.Decode(&resp)
	if err != nil {
		t.Errorf("expected err to be nil got %v", err)
	}

	if resp.Key != "key1" {
		t.Errorf("expected key1 but got %v", resp.Key)
	}
	if resp.Value != "val1" {
		t.Errorf("expected val1 but got %v", resp.Value)
	}

	resp2, err := getMemoryRequest("key1")
	if err != nil {
		t.Errorf("expected err to be nil got %v", err)
	}
	if resp2.Value != "val1" {
		t.Errorf("expected val1 but got %v", resp2.Value)
	}

	resp3, err := getMemoryRequest("key2")
	if resp3.Value != "" {
		t.Errorf("expected empty string but got %v", resp3.Value)
	}

}
