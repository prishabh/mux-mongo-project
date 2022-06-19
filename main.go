package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prishabh/mux-mongo-project/libs/database/driver/mongodb"
	"github.com/prishabh/mux-mongo-project/models"
	"github.com/prishabh/mux-mongo-project/responses"
	"go.mongodb.org/mongo-driver/bson"
	"gopkg.in/validator.v2"
)

var (
	mongodbClient *mongodb.Client
)

func handleCreate(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("content-type", "application/json")
	var data []models.Data
	//validate the request body
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		response := responses.Response{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    map[string]interface{}{"data": err.Error()},
		}
		json.NewEncoder(rw).Encode(response)
		return
	}
	//use the validator library to validate required fields
	if validationErr := validator.Validate(data); validationErr != nil {
		rw.WriteHeader(http.StatusBadRequest)
		response := responses.Response{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    map[string]interface{}{"data": validationErr.Error()},
		}
		json.NewEncoder(rw).Encode(response)
		return
	}
	for _, d := range data {
		_, err := mongodbClient.Insert(d)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.Response{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data:    map[string]interface{}{"data": err.Error()},
			}
			json.NewEncoder(rw).Encode(response)
			return
		}
	}
	rw.WriteHeader(http.StatusCreated)
	response := responses.Response{
		Status:  http.StatusCreated,
		Message: "success",
	}
	json.NewEncoder(rw).Encode(response)
}

func handleRead(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("content-type", "application/json")
	params := r.URL.Query()
	first_name := params.Get("first_name")
	last_name := params.Get("last_name")
	car_manufactur := params.Get("car_manufactur")
	city := params.Get("city")
	var filter interface{}
	if len(first_name) > 0 && len(last_name) > 0 {
		filter = bson.D{{"first_name", first_name}, {"last_name", last_name}}
	} else if len(car_manufactur) > 0 {
		filter = bson.D{{"car_manufactur", car_manufactur}}
	} else if len(city) > 0 {
		if _, err := mongodbClient.CreateIndex("address", "text"); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.Response{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data:    map[string]interface{}{"data": err.Error()},
			}
			json.NewEncoder(rw).Encode(response)
			return
		}
		filter = bson.D{{"$text", bson.D{{"$search", city}}}}
	} else {
		filter = bson.D{}
	}
	var results []models.Data
	if err := mongodbClient.Query(filter, &results); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		response := responses.Response{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    map[string]interface{}{"data": err.Error()},
		}
		json.NewEncoder(rw).Encode(response)
		return
	}
	rw.WriteHeader(http.StatusOK)
	response := responses.Response{
		Status:  http.StatusOK,
		Message: "success",
		Data:    map[string]interface{}{"data": results},
	}
	json.NewEncoder(rw).Encode(response)
}

func initializeRouter() {
	log.Print("Initializing routes")
	r := mux.NewRouter()
	r.HandleFunc("/", handleCreate).Methods("POST")
	r.HandleFunc("/", handleRead).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func main() {
	log.Print("Starting go webserver")
	var err error
	// mongodb config
	mongodbConfig := new(mongodb.Config)
	mongodbConfig.Address = "localhost:27017"
	mongodbConfig.User = "root"
	mongodbConfig.Password = "example"
	mongodbConfig.Database = "test"
	mongodbConfig.Collection = "data"
	// get mongodb client
	mongodbClient, err = mongodb.NewClient(mongodbConfig)
	if err != nil {
		log.Fatal(err)
	}
	initializeRouter()
}
