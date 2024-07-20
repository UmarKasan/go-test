package main

import (
	"context"
	"encoding/json"
	"go-test2/helper"
	"go-test2/models"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var collection *mongo.Collection

func main() {
	// Connect to MongoDB
	collection = helper.ConnectDB()

	//Init Router
	route := mux.NewRouter()
	// For CORS
	route.Use(mux.CORSMethodMiddleware(route))

	route.HandleFunc("/api/books", func(w http.ResponseWriter, route *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		getBooks(w, route)
	}).Methods("GET", "OPTIONS")
	// End CORS

	// routes
	route.HandleFunc("/api/books", getBooks).Methods("GET")
	route.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	route.HandleFunc("/api/books", createBook).Methods("POST")
	route.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	route.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	// set port address
	log.Fatal(http.ListenAndServe(":8000", route))

}

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Create a slice to hold the books
	var books []models.Book

	// Find all books
	cur, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		helper.GetError(err, w)
		return
	}
	defer cur.Close(context.TODO())

	// Iterate through the cursor and decode each document
	for cur.Next(context.TODO()) {
		var book models.Book
		err := cur.Decode(&book)
		if err != nil {
			log.Println(err)
			continue
		}
		books = append(books, book)
	}

	if err := cur.Err(); err != nil {
		helper.GetError(err, w)
		return
	}

	// Encode and return the books
	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	// set header.
	w.Header().Set("Content-Type", "application/json")

	var book models.Book
	// we get params with mux.
	var params = mux.Vars(r)

	// string to primitive.ObjectID
	id, _ := primitive.ObjectIDFromHex(params["id"])

	// We create filter. If it is unnecessary to sort data for you, you can use bson.M{}
	filter := bson.M{"_id": id}
	err := collection.FindOne(context.TODO(), filter).Decode(&book)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(book)
}

func createBook(w http.ResponseWriter, route *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var book models.Book

	// we decode our body request params
	_ = json.NewDecoder(route.Body).Decode(&book)

	// insert our book model.
	result, err := collection.InsertOne(context.TODO(), book)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(result)
}
func updateBook(w http.ResponseWriter, route *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(route)

	//Get id from parameters
	id, _ := primitive.ObjectIDFromHex(params["id"])

	var book models.Book

	// Create filter
	filter := bson.M{"_id": id}

	// Read update model from body request
	_ = json.NewDecoder(route.Body).Decode(&book)

	// prepare update model.
	update := bson.D{
		{"$set", bson.D{
			{"isbn", book.Isbn},
			{"title", book.Title},
			{"author", bson.D{
				{"firstname", book.Author.FirstName},
				{"lastname", book.Author.LastName},
			}},
		}},
	}

	err := collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&book)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	book.ID = id

	json.NewEncoder(w).Encode(book)
}
func deleteBook(w http.ResponseWriter, route *http.Request) {
	// Set header
	w.Header().Set("Content-Type", "application/json")

	// get params
	var params = mux.Vars(route)

	// string to primitve.ObjectID
	id, err := primitive.ObjectIDFromHex(params["id"])

	// prepare filter.
	filter := bson.M{"_id": id}

	deleteResult, err := collection.DeleteOne(context.TODO(), filter)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(deleteResult)
}
