package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
) 

// Create a struct that has a log object 
type httpServer struct {
	 Log *Log
}

// Create a function that returns a reference to that log
func newHTTPServer() *httpServer {
	// Return a ptr to a new struct literal
	return &httpServer{ Log: NewLog()}
}

// Create structs for json unmarshaling and marshaling
type ProduceRequest struct {
	Record Record `json:record`
}

type ProduceResponse struct {
	Offset uint64 `json:"offset"`
}

type ConsumeRequest struct {
	Offset uint64 `json:"offset"`
}

type ConsumeResponse struct {
	Record Record `json:record`
}

func (server *httpServer) handleProduce(writer http.ResponseWriter, reader *http.Request) {
	// Decode the json and marshall it onto our ProduceRequest struct and do error handling
	var request ProduceRequest
	err := json.NewDecoder(reader.Body).Decode(&request)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	// Try to append the log. If we can't append it, return a 500
	offset, err := server.Log.Append(request.Record)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	// Create the response object and add the output to the http writer. If we can't marshal the JSON, return a 500
	response := ProduceResponse{Offset: offset}
	err = json.NewEncoder(writer).Encode(response)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (server *httpServer) handleConsume(writer http.ResponseWriter, reader *http.Request) {
	// Decode the json and marshall it onto our ProduceRequest struct and do error handling
	var request ConsumeRequest
	err := json.NewDecoder(reader.Body).Decode(&request)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	// Try to read the log. If we can't find it, return a 404
	record, err := server.Log.Read(request.Offset)
	if err == ErrOffsetNotFound {
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	response := ConsumeResponse{Record: record}
	err = json.NewEncoder(writer).Encode(response)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func NewHttpServer(address string) *http.Server {
	// Get a new http server struct
	httpServer := newHTTPServer()
	// Get a new router via mux lib
	router := mux.NewRouter()

	// Register routes and handler functions
	router.HandleFunc("/api/v1/log", httpServer.handleProduce).Methods("POST")
	router.HandleFunc("/api/v1/log", httpServer.handleConsume).Methods("GET")

	// Return our cooked up sever
	return &http.Server{
		Addr: address,
		Handler: router,
	}
}