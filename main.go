package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gorilla/mux"
)

// CustomError is a custom error type
type CustomError struct {
	Message string
}

// Error implements the error interface for CustomError
func (e *CustomError) Error() string {
	return e.Message
}

// Response is a struct to represent the JSON response
type Response struct {
	Message string `json:"message"`
}

//dd:span
func apiHandler(w http.ResponseWriter, r *http.Request) {

	// Set the content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Create a Response struct
	response := Response{Message: "Hello, this is a simple GET API!"}
	log.Println("Simple get API Success!!")

	// Encode the Response struct to JSON and write it to the response writer
	json.NewEncoder(w).Encode(response)
}

//dd:span
func getErrorRequestHandler(w http.ResponseWriter, r *http.Request) {

	// Set the content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Return a custom error
	customError := &CustomError{Message: "Custom Error: Something went wrong"}
	w.WriteHeader(http.StatusInternalServerError)

	// Log the error message
	log.Printf("Error Request API triggered with message: %s\n", customError.Message)
	debug.PrintStack() // Capture and print the stack trace

	// Encode the custom error message to JSON and write it to the response writer
	json.NewEncoder(w).Encode(Response{Message: customError.Message})
}

func main() {
	router := mux.NewRouter()

	// Register the API handler function for the "/api" route using gorilla/mux
	router.HandleFunc("/api", apiHandler).Methods("GET")

	// Register the handler function for the "/getErrorRequest" route
	router.HandleFunc("/getErrorRequest", getErrorRequestHandler).Methods("GET")

	// Attach the router to the default serve mux
	http.Handle("/", router)

	// Start the server on port 8080
	fmt.Println("Server is running on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("Error:", err)
	}
}
