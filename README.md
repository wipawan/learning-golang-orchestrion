# learning-golang-orchestrion

- [Orchestrion](https://github.com/DataDog/orchestrion) is in private beta with Datadog. CAA 9th March 2024.
- For more in depth code sample refer [here](https://github.com/DataDog/go-sample-app).
- There are 2 ways you can use orchestrion to help with instrumentation:
1. Locally, install orchestrion on machine and run the tool
2. At build time.

## Intro
- This code base shows how you can make use of orchestrion on the Dockerfile layer to help auto instrument Golang at **build time**. 
- Orchestrion supports some library for [Auto Instrumentation](https://github.com/DataDog/orchestrion?tab=readme-ov-file#supported-libraries).
- You are still expected to go through your golang codes to annotate your codes for this to work.
  - For libraries that are not in the list, you will be able to annotate your code files with //dd:span my:tag where it represents //dd:span <custom span tag>. 

## Inspect main.go
- The outcome of the orchestrion instrumentation is that it should have 3 spans in the flamegraph. 
- All we did was to annotate //dd:span my:tag on func apiHandler and getRequestHandler. This helps us achieve the instrumentation for that function.
- Since we are using the net/http and gorilla/mux library that orchestrion supports for [automatic instrumentation](https://github.com/DataDog/orchestrion?tab=readme-ov-file#supported-libraries), we do not need to annotate those.

## Inspect afterorchestrion.go.example
- This is how the main.go file will look like after orchestrion ["automagically"](https://github.com/DataDog/orchestrion?tab=readme-ov-file#how-it-works) instruments your code.

## How it looks like in Datadog APM FlameGraph after orchestrion
![Orchestrion Scenario 2](https://github.com/jon94/learning-golang-orchestrion/assets/40360784/c4498456-8c8f-40df-811d-7b85a33da33c)

## See it in action
1. Clone the repo
```
git clone https://github.com/jon94/learning-golang-orchestrion.git
```
2. Replace the Datadog API Key in docker-compose.yaml
![image](https://github.com/jon94/learning-golang-orchestrion/assets/40360784/ec42f1fb-bcc6-4e23-bfeb-11fcf7eb4b86)
3. Set ENV Variable [DD_SITE](https://docs.datadoghq.com/agent/troubleshooting/site/) if required (depending on your data centre with Datadog)
![image](https://github.com/jon94/learning-golang-orchestrion/assets/40360784/44a0c8fe-29cf-473a-98a5-441a03737e31)
4. Run docker compose
```
docker compose up -d --force-recreate --no-deps --build
```
5. Generate traffic by hitting curl -v http://localhost:5000/apiRequest and curl -v http://localhost:5000/getErrorRequest
6. After you are done, head to the Datadog APM Platform to see the Traces.
```
docker compose down
```
## For Inline Functions
Depending on how your Golang code is structured. You might need to use a combination of Orchestrion and [custom instrumentation](https://docs.datadoghq.com/tracing/trace_collection/custom_instrumentation/dd_libraries/go/#adding-spans)

<details>
<summary>Click to toggle for more info</summary>

<details>
<summary>Example Code Block</summary>

- In this case, apiHandler and getRequestHandler functions are not written outside the func main() block. They are written as inline functions inside func main(). Orchestrion is not yet able to catch those.
- We therefore, created 2 manual spans for the 2 functions.
  
```Go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/gorilla/mux"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
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
func main() {
	// Create a new router using gorilla/mux
	router := mux.NewRouter()

	// Define a handler function for the API endpoint
	apiHandler := func(w http.ResponseWriter, r *http.Request) {
		// Start a trace span for this handler
		ctx := r.Context()
		span, ctx := tracer.StartSpanFromContext(ctx, "apiHandler", tracer.ResourceName("/simplegetapi"))
		defer span.Finish()

		// Set the content type to JSON
		w.Header().Set("Content-Type", "application/json")

		// Create a Response struct
		response := Response{Message: "Hello, this is a simple GET API!"}
		log.Println("Simple get API Success!!")

		// Encode the Response struct to JSON and write it to the response writer
		json.NewEncoder(w).Encode(response)
	}

	// Define a handler function for the /getErrorRequest endpoint
	getErrorRequestHandler := func(w http.ResponseWriter, r *http.Request) {
		// Start a trace span for this handler
		ctx := r.Context()
		span, ctx := tracer.StartSpanFromContext(ctx, "getErrorRequestHandler", tracer.ResourceName("/getErrorRequest"))
		defer span.Finish()

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
```
</details>
</details>

## Credits
- Sin Ta: For debugging the Dockerfile with me to make it work.
- Sho Uchida: For bumping Orchestrion for Golang Auto Instrumentation.
