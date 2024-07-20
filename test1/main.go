package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MarketplaceItem struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName          string             `bson:"FirstName" json:"FirstName"`
	LastName           string             `bson:"LastName" json:"LastName"`
	Product            string             `bson:"Product" json:"Product"`
	Quantity           int                `bson:"Quantity" json:"Quantity"`
	Condition          string             `bson:"Condition" json:"Condition"`
	CollectionLocation string             `bson:"Collection-Location" json:"CollectionLocation"`
}

var client *mongo.Client
var collection *mongo.Collection

func main() {
	// Set up MongoDB connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")

	// Get a handle for your collection
	collection = client.Database("SCRAP").Collection("Marketplace")

	// Ensure the collection exists and add some sample data
	ensureCollectionAndData(ctx)

	// Set up router
	r := mux.NewRouter()
	r.HandleFunc("/api/items", getItems).Methods("GET")
	r.HandleFunc("/api/items", createItem).Methods("POST")
	r.HandleFunc("/api/items/{id}", getItem).Methods("GET")
	r.HandleFunc("/api/items/{id}", updateItem).Methods("PUT")
	r.HandleFunc("/api/items/{id}", deleteItem).Methods("DELETE")

	log.Println("Server is starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func ensureCollectionAndData(ctx context.Context) {
	// Check if the collection is empty
	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	// If the collection is empty, insert some sample data
	if count == 0 {
		sampleItems := []interface{}{
			MarketplaceItem{
				FirstName:          "Mat",
				LastName:           "Motor",
				Product:            "Energizer Battery",
				Quantity:           15,
				Condition:          "Good",
				CollectionLocation: "123 Tampines Street 33",
			},
			MarketplaceItem{
				FirstName:          "Jane",
				LastName:           "Doe",
				Product:            "Aluminum Cans",
				Quantity:           50,
				Condition:          "Used",
				CollectionLocation: "456 Orchard Road",
			},
		}

		_, err := collection.InsertMany(ctx, sampleItems)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Inserted sample data into the collection")
	}
}

func getItems(w http.ResponseWriter, r *http.Request) {
	var items []MarketplaceItem
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var item MarketplaceItem
		cursor.Decode(&item)
		items = append(items, item)
	}

	json.NewEncoder(w).Encode(items)
}

func createItem(w http.ResponseWriter, r *http.Request) {
	var item MarketplaceItem
	json.NewDecoder(r.Body).Decode(&item)

	result, err := collection.InsertOne(context.Background(), item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	item.ID = result.InsertedID.(primitive.ObjectID)
	json.NewEncoder(w).Encode(item)
}

func getItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])

	var item MarketplaceItem
	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(item)
}

func updateItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])

	var item MarketplaceItem
	json.NewDecoder(r.Body).Decode(&item)

	_, err := collection.UpdateOne(
		context.Background(),
		bson.M{"_id": id},
		bson.D{
			{"$set", bson.D{
				{"FirstName", item.FirstName},
				{"LastName", item.LastName},
				{"Product", item.Product},
				{"Quantity", item.Quantity},
				{"Condition", item.Condition},
				{"Collection-Location", item.CollectionLocation},
			}},
		},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	item.ID = id
	json.NewEncoder(w).Encode(item)
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])

	_, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
