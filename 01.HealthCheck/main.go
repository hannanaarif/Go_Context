// package main

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"time"

// 	"github.com/go-chi/chi/v5"
// 	"github.com/go-chi/chi/v5/middleware"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// func connectToMongoDB(ctx context.Context) (*mongo.Client, error) {
// 	clientOpts := options.Client().
// 		ApplyURI("mongodb+srv://hannanaarif:India%402025@mongogo.dvnres0.mongodb.net/?retryWrites=true&w=majority&appName=MongoGO").SetServerSelectionTimeout(5 * time.Second)
// 	// Implement connection retry logic 
// 	maxRetries := 3
// 	var client *mongo.Client
// 	var err error

// 	for i := 0; i < maxRetries; i++ {
// 		client, err = mongo.Connect(ctx, clientOpts)
// 		if err == nil {
// 			// Try to ping the server to verify the connection
// 			if err = client.Ping(ctx, nil); err == nil {
// 				log.Println("Successfully connected to MongoDB")
// 				return client, nil
// 			}
// 		}
// 		log.Printf("Failed to connect to MongoDB (attempt %d/%d): %v", i+1, maxRetries, err)
// 		if i < maxRetries-1 {
// 			time.Sleep(2 * time.Second) // Wait before retrying
// 		}
// 	}
// 	return nil, fmt.Errorf("failed to connect to MongoDB after %d attempts: %v", maxRetries, err)
// }

// func main() {
// 	r := chi.NewRouter()
// 	r.Use(middleware.Logger)

// 	// Create a base context with a longer timeout for initial connection
// 	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
// 	defer cancel()

// 	// Connect to MongoDB with retry logic
// 	client, err := connectToMongoDB(ctx)
// 	if err != nil {
// 		log.Fatalf("Failed to connect to MongoDB: %v", err)
// 	}
// 	defer func() {
// 		if err := client.Disconnect(context.Background()); err != nil {
// 			log.Printf("Error disconnecting from MongoDB: %v", err)
// 		}
// 	}()

// 	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
// 		w.Write([]byte("Welcome"))
// 	})

// 	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
// 		healthCtx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
// 		defer cancel()

// 		err := client.Ping(healthCtx, nil)
// 		if err != nil {
// 			log.Printf("Health check failed: %v", err)
// 			http.Error(w, fmt.Sprintf("MongoDB connection failed: %v", err), http.StatusServiceUnavailable)
// 			return
// 		}
// 		w.Write([]byte("MongoDB connection healthy"))
// 	})

// 	fmt.Println("Server is running on http://localhost:3000")
// 	if err := http.ListenAndServe(":3000", r); err != nil {
// 		log.Fatal(err)
// 	}
// }

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func connectToMongoDB() (*mongo.Client, error) {
	uri := "mongodb+srv://hannanaarif:India%402025@mongogo.dvnres0.mongodb.net/?retryWrites=true&w=majority"
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// Connect with a timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	// Ping the server to verify the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	log.Println("Successfully connected to MongoDB")
	return client, nil
}

func main() {
	// Connect to MongoDB
	client, err := connectToMongoDB()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	// Set up a simple HTTP server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome"))
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		// Check MongoDB health
		err := client.Ping(context.Background(), nil)
		if err != nil {
			log.Printf("Health check failed: %v", err)
			http.Error(w, "MongoDB connection failed", http.StatusServiceUnavailable)
			return
		}
		w.Write([]byte("MongoDB connection healthy"))
	})

	// Start the server
	log.Println("Server is running on http://localhost:3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}
