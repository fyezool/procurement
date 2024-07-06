package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Procurement struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	Status      string             `bson:"status" json:"status"`
	DateCreated time.Time          `bson:"date_created" json:"date_created"`
	DateUpdated time.Time          `bson:"date_updated" json:"date_updated"`
}

func getProcurements(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var procurements []Procurement
	collection := client.Database("procurementdb").Collection("procurements")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var procurement Procurement
		cursor.Decode(&procurement)
		procurements = append(procurements, procurement)
	}
	if err := cursor.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(w).Encode(procurements)
}

func createProcurement(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var procurement Procurement
	json.NewDecoder(r.Body).Decode(&procurement)
	procurement.DateCreated = time.Now()
	procurement.DateUpdated = time.Now()
	collection := client.Database("procurementdb").Collection("procurements")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := collection.InsertOne(ctx, procurement)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(w).Encode(result)
}
