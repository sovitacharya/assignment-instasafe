package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestPost(t *testing.T) {
	var testReq TransactionReq
	sampleErrorReq := `{"timestamp":"abcdefgh","amount":"efg"}`

	req, err := http.NewRequest("POST", "/transactions", bytes.NewBuffer([]byte(sampleErrorReq)))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Transactions)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != 422 {
		t.Error("handler returned wrong status code: got ", status, "want", 422)
	}
	testReq.Amount = 10
	presentTime := time.Now().UTC()
	testReq.Timestamp = presentTime
	// reqBody, _ := json.Marshal(testReq)
	reqBody, _ := json.Marshal(testReq)
	req, err = http.NewRequest("POST", "/transactions", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != 201 {
		t.Error("handler returned wrong status code: got ", status, "want", 201)
	}
	testReq.Amount = 100
	nextTime := presentTime.Add(100 * time.Second)
	testReq.Timestamp = nextTime
	reqBody, _ = json.Marshal(testReq)
	req, err = http.NewRequest("POST", "/transactions", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != 422 {
		t.Error("handler returned wrong status code: got ", status, "want", 422)
	}
	testReq.Amount = 100
	nextTime = presentTime.Add(-60 * time.Second)
	testReq.Timestamp = nextTime
	reqBody, _ = json.Marshal(testReq)
	req, err = http.NewRequest("POST", "/transactions", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != 204 {
		t.Error("handler returned wrong status code: got ", status, "want", 204)
	}
}
func TestGet(t *testing.T) {
	req, err := http.NewRequest("GET", "/statistics", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Statistics)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Error("handler returned wrong status code: got ", status, "want", http.StatusOK)
	}

}

func TestDelete(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/transactions", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Statistics)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Error("handler returned wrong status code: got ", status, "want", http.StatusOK)
	}

}
