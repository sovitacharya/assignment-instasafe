package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.elastic.co/apm/module/apmhttp"
	"gopkg.in/go-playground/validator.v9"
)

type TransactionReq struct {
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
}
type Response struct {
	Sumvalue   float64 `json:"sum"`
	Average    float64 `json:"avg"`
	Maxvalue   float64 `json:"max"`
	Minvalue   float64 `json:"min"`
	Countvalue float64 `json:"count"`
}
type SliceAmount struct {
	TimeStamp time.Time
	Resvalue  Response
}

var Value []TransactionReq
var ResponseSlice []SliceAmount

var Input TransactionReq
var responsesetStruct Response
var dataSlicestruct SliceAmount

func main() {
	mux := NewRouter()
	http.Handle("/", mux)
	http.ListenAndServe(":8080", apmhttp.Wrap(mux))
}

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", appstart).Methods("POST")
	r.HandleFunc("/transactions", Transactions).Methods("POST")
	r.HandleFunc("/statistics", Statistics).Methods("GET")
	r.HandleFunc("/transactions", DeleteTransactions).Methods("DELETE")
	return r
}

func appstart(res http.ResponseWriter, req *http.Request) {
	_, err := io.WriteString(res, "Hi!, This is the assignment application.")
	log.Println("Successfully Listening on 8080 port ")
	if err != nil {
		return
	}
}
func Transactions(res http.ResponseWriter, req *http.Request) {
	// var Input TransactionReq
	// var responsesetStruct Response
	// var dataSlicestruct SliceAmount
	err := json.NewDecoder(req.Body).Decode(&Input)
	if err != nil {
		fmt.Println("Incorrect Json Req", err.Error())
		res.WriteHeader(422)
		json.NewEncoder(res).Encode("Incorrect Json Req")
		return
	}

	Valid := TransactionReq{
		Amount:    Input.Amount,
		Timestamp: Input.Timestamp,
	}
	v := validator.New()
	err = v.Struct(Valid)
	if err != nil {
		fmt.Println("Invalid Json Req", err.Error())
		res.WriteHeader(400)
		json.NewEncoder(res).Encode("Invalid Json Req")
		return
	}
	fmt.Println("req body", Input)
	presentTime := time.Now().UTC()
	if presentTime.UTC().Sub(Input.Timestamp).Seconds() < 0 {
		fmt.Println("future Transaction")
		res.WriteHeader(422)
		json.NewEncoder(res).Encode("Future Transaction")
		return
	}
	if presentTime.Sub(Input.Timestamp).Seconds() > 60 {
		fmt.Println("60 sec oldTransaction")

		Value = append(Value, Input)
		fmt.Println("valuecheck", Value)
		res.WriteHeader(204)
		return
	}
	Value = append(Value, Input)
	var data TransactionReq
	var max float64
	var sum float64
	var min float64
	var count float64
	fmt.Println("value", Value)
	dataCount := len(Value)
	fmt.Println("datacount value", dataCount)
	if dataCount != 0 {
		for i := 0; i < dataCount; i++ {
			data = Value[i]
			if time.Now().UTC().Sub(data.Timestamp).Seconds() <= 60 {
				sum = sum + data.Amount
				fmt.Println("sum", sum)
				count++
				fmt.Println("count", count)
				if data.Amount > max {
					max = data.Amount
					fmt.Println("max", max)
				}
				if min == 0 {
					min = data.Amount
					fmt.Println("min", min)
				} else if data.Amount < min {
					min = data.Amount
				}

			}
		}
	}

	responsesetStruct.Maxvalue = max
	responsesetStruct.Minvalue = min
	responsesetStruct.Sumvalue = sum
	responsesetStruct.Countvalue = count
	if count != 0 {
		responsesetStruct.Average = sum / count
	}
	dataSlicestruct.TimeStamp = Input.Timestamp
	dataSlicestruct.Resvalue = responsesetStruct
	ResponseSlice = append(ResponseSlice, dataSlicestruct)
	fmt.Println("Data slice ", ResponseSlice)
	res.WriteHeader(201)
	json.NewEncoder(res).Encode("Successfuly Saved")
	return
}

func Statistics(res http.ResponseWriter, req *http.Request) {
	var response Response
	fmt.Println("response slice", ResponseSlice)
	fmt.Println("Get req started")
	presentTIme := time.Now().UTC()
	sliceLength := len(ResponseSlice)

	if sliceLength != 0 {
		fmt.Println("entry into 1st loop")
		reqSlice := ResponseSlice[sliceLength-1]
		if presentTIme.Sub(reqSlice.TimeStamp).Seconds() <= 60 {
			fmt.Println("enter into second loop")
			response = reqSlice.Resvalue
			res.WriteHeader(200)
			fmt.Println("response for get req", response)
			json.NewEncoder(res).Encode(response)
			return
		} else {
			res.WriteHeader(200)
			json.NewEncoder(res).Encode(response)
			return
		}
	}
	res.WriteHeader(200)
	json.NewEncoder(res).Encode(response)
	return
}

func DeleteTransactions(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Delete all req stored")
	fmt.Println("Data stored are ", Value)
	Value = nil
	fmt.Println("data after clearing", Value)
	res.WriteHeader(204)

}
