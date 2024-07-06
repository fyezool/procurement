package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Quote struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	ProcurementID primitive.ObjectID `bson:"procurement_id" json:"procurement_id"`
	SupplierID    primitive.ObjectID `bson:"supplier_id" json:"supplier_id"`
	Items         []QuoteItem        `bson:"items" json:"items"`
	TotalCost     float64            `bson:"total_cost" json:"total_cost"`
}

type QuoteItem struct {
	ItemID    primitive.ObjectID `bson:"item_id" json:"item_id"`
	Quantity  int                `bson:"quantity" json:"quantity"`
	UnitPrice float64            `bson:"unit_price" json:"unit_price"`
}

func getQuotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var quotes []Quote
	collection := client.Database("procurementdb").Collection("quotes")
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
		var quote Quote
		cursor.Decode(&quote)
		quotes = append(quotes, quote)
	}
	if err := cursor.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(w).Encode(quotes)
}

func createQuote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var quote Quote
	json.NewDecoder(r.Body).Decode(&quote)
	collection := client.Database("procurementdb").Collection("quotes")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := collection.InsertOne(ctx, quote)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(w).Encode(result)
}
