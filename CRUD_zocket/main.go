package main

import (
        "context"
        "go.mongodb.org/mongo-driver/mongo"
        "go.mongodb.org/mongo-driver/bson" 
        "go.mongodb.org/mongo-driver/mongo/options"
        // "go.mongodb.org/mongo-driver/mongo/readpref"
        "fmt"
        "os"
        "github.com/gorilla/mux"
        _ "github.com/joho/godotenv/autoload"
        "encoding/json"
        "log"
        "net/http"
        // "time"

)

type Item struct {
    ID        string    `json:"id,omitempty" bson:"_id,omitempty"`
    BooktName string    `json:"bookName,omitempty" bson:"bookName,omitempty"`
    AuthorName  string    `json:"authorName,omitempty" bson:"authorName,omitempty"`
    // CreatedAt time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
}

var collection *mongo.Collection
func main() {
        // connect to the MongoDB server
        var url = os.Getenv("MONGO_URL")
        client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(url))
        if err != nil {
                panic(err)
        }
        if err := client.Ping(context.Background(), nil); err != nil {
            panic(err)
        }

        // get a handle for your collection
        collection = client.Database("Item").Collection("items")

        router := mux.NewRouter()
		router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Welcome to the HomePage!")
		})

        // Hanclers for different CRUD operations
        router.HandleFunc("/createitem", createItem).Methods("POST")
        router.HandleFunc("/getitem", getItems).Methods("GET")
        router.HandleFunc("/getitem/{id}", getItem).Methods("GET")
        router.HandleFunc("/updateitem/{id}", updateItem).Methods("PUT")
        router.HandleFunc("/deleteitem/{id}", deleteItem).Methods("DELETE")
        fmt.Println("Server is running on port 8888")
		log.Fatal(http.ListenAndServe(":8888", router))
        
}   
// CreatItem function
func createItem(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    var item Item
    json.NewDecoder(r.Body).Decode(&item)
    result, err := collection.InsertOne(context.Background(), item)
    if err != nil {
        log.Fatal(err)
    }
    json.NewEncoder(w).Encode(result)
}
func getItems(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    var items []Item
    cur, err := collection.Find(context.Background(), bson.M{})
    if err != nil {
        log.Fatal(err)
    }
    defer cur.Close(context.Background())
    for cur.Next(context.Background()) {
        var item Item
        cur.Decode(&item)
        items = append(items, item)
    }
    if err := cur.Err(); err != nil {
        log.Fatal(err)
    }
    json.NewEncoder(w).Encode(items)
}

func getItem(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    params := mux.Vars(r)
    id := params["id"]
    var item Item
    err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&item)
    if err != nil {
        log.Fatal(err)
    }
    json.NewEncoder(w).Encode(item)
}

func updateItem(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    params := mux.Vars(r)
    id := params["id"]
    var item Item
    json.NewDecoder(r.Body).Decode(&item)
    result, err := collection.ReplaceOne(context.Background(), bson.M{"_id": id}, item)
    if err != nil {
        log.Fatal(err)
    }
    json.NewEncoder(w).Encode(result)
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    params := mux.Vars(r)
    id := params["id"]
    result, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})
    if err != nil {
        log.Fatal(err)
    }
    json.NewEncoder(w).Encode(result)
}

