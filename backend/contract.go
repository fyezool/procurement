package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Contract struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	ProcurementID primitive.ObjectID `bson:"procurement_id" json:"procurement_id"`
	SupplierID    primitive.ObjectID `bson:"supplier_id" json:"supplier_id"`
	Items         []ContractItem     `bson:"items" json:"items"`
	TotalCost     float64            `bson:"total_cost" json:"total_cost"`
}

type ContractItem struct {
	ItemID    primitive.ObjectID `bson:"item_id" json:"item_id"`
	Quantity  int                `bson:"quantity" json:"quantity"`
	UnitPrice float64            `bson:"unit_price" json:"unit_price"`
}

func getContracts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var contracts []Contract
	collection := client.Database("procurementdb").Collection("contracts")
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
		var contract Contract
		cursor.Decode(&contract)
		contracts = append(contracts, contract)
	}
	if err := cursor.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(w).Encode(contracts)
}

func createContract(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var contract Contract
	json.NewDecoder(r.Body).Decode(&contract)
	collection := client.Database("procurementdb").Collection("contracts")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := collection.InsertOne(ctx, contract)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(w).Encode(result)
}
