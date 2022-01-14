package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>It works!</h1>"))
}

func mongoApi(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		var req MongoRequest
		err := decoder.Decode(&req)
		if err != nil {
			result := &CustomError{Error: "invalid request"}
			sendJson(&w, result)
		}
		records, err := getRecords(req.StartDate, req.EndDate, req.MinCount, req.MaxCount)
		if err != nil {
			result := &CustomError{Error: err.Error()}
			sendJson(&w, result)
		}
		sendJson(&w, &MongoResponse{Records: &records, Code: 0, Msg: "Success"})
	} else {
		w.Write([]byte("<h1>!</h1>"))
	}
}

func memoryApi(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		key := r.URL.Query().Get("key")
		res := getMemoryData(key)
		if res != "" {
			result := &InMemory{Key: key, Value: res}
			sendJson(&w, result)
		} else {
			result := &CustomError{Error: "key not found!"}
			sendJson(&w, result)
		}
	} else if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		var req InMemory
		err := decoder.Decode(&req)
		if err != nil {
			result := &CustomError{Error: "invalid request"}
			sendJson(&w, result)
		}
		setMemoryData(req.Key, req.Value)
		sendJson(&w, req)
	} else {
		w.Write([]byte("<h1>!</h1>"))
	}
}

func sendJson(w *http.ResponseWriter, any interface{}) {
	(*w).Header().Set("Content-Type", "application/json")
	bytes, _ := json.Marshal(any)
	fmt.Fprintf((*w), string(bytes))
}
