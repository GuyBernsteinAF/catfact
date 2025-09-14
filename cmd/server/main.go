package main

import (
	"bytes"
	"catfacts/docs"
	"catfacts/internal"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type SuccessResponse struct {
	Message string   `json:"message"`
	Facts   []string `json:"facts"`
}

func captureStdout(f func()) string {
	// Save the original stdout
	originalStdout := os.Stdout

	// Create a pipe to capture output
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run the function in a goroutine to prevent blocking
	done := make(chan bool)
	var buf bytes.Buffer

	go func() {
		io.Copy(&buf, r)
		done <- true
	}()

	// Call the function that writes to stdout
	f()

	// Restore original stdout and close the writer
	w.Close()
	os.Stdout = originalStdout

	// Wait for the capture to complete
	<-done

	return buf.String()
}

func phaseOneAPI(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, captureStdout(internal.PhaseOne))
}

func phaseTwoAPI(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, captureStdout(internal.PhaseTwo))
}

func phaseThreeAPI(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, captureStdout(internal.PhaseThree))
}

func phaseFourAPI(w http.ResponseWriter, req *http.Request) {
	name := req.URL.Query().Get("name")
	amount := req.URL.Query().Get("amount")
	w.Header().Set("Content-Type", "application/json")

	if amount == "" {
		amount = "1"
	}

	if !validate(w, amount, name) {
		return
	}

	res := SuccessResponse{Message: "Hello " + name + ", here are you cat facts"}
	intAmount, _ := strconv.Atoi(amount)
	res.Facts = internal.PhaseFour(intAmount)
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(res)
	return

}

func validate(w http.ResponseWriter, amount string, name string) bool {
	if am, err := strconv.Atoi(amount); err != nil || am <= 0 || am > 10 {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(ErrorResponse{"amount must be an integer between 1 and 10 (or not required)"})
		return false
	}

	if name == "" {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(ErrorResponse{"name is required"})
		return false
	}

	if len(name) > 32 || strings.Contains(name, " ") {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(ErrorResponse{"name is a single word with length 1-32"})
		return false
	}

	for _, c := range name {
		if !unicode.IsLetter(rune(c)) {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(ErrorResponse{"name should be alphabetic"})
			return false
		}
	}
	return true
}

func headers(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

// Serve the Swagger spec
func swaggerSpecHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/x-yaml")
	fmt.Fprint(w, docs.SwaggerSpec)
}

// Serve Swagger UI
func swaggerUIHandler(w http.ResponseWriter, req *http.Request) {
	html := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Cat Facts API - Swagger UI</title>
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5.9.0/swagger-ui.css">
    <style>
        body {
            margin: 0;
            padding: 0;
        }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5.9.0/swagger-ui-bundle.js"></script>
    <script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5.9.0/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function() {
            const ui = SwaggerUIBundle({
                url: "/swagger/spec",
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout"
            });
            window.ui = ui;
        }
    </script>
</body>
</html>
`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

// Home page with API documentation links
func homeHandler(w http.ResponseWriter, req *http.Request) {
	html := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Cat Facts API</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
            max-width: 800px;
            margin: 50px auto;
            padding: 20px;
            background-color: #f5f5f5;
        }
        h1 {
            color: #333;
        }
        .card {
            background: white;
            border-radius: 8px;
            padding: 20px;
            margin: 20px 0;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        a {
            color: #007bff;
            text-decoration: none;
        }
        a:hover {
            text-decoration: underline;
        }
        ul {
            line-height: 1.8;
        }
        .swagger-link {
            display: inline-block;
            background: #85ea2d;
            color: #173647;
            padding: 10px 20px;
            border-radius: 4px;
            margin-top: 10px;
            font-weight: bold;
        }
        .swagger-link:hover {
            background: #75da1d;
            text-decoration: none;
        }
    </style>
</head>
<body>
    <h1>🐱 Cat Facts API</h1>
    
    <div class="card">
        <h2>Welcome!</h2>
        <p>This API provides random cat facts through different retrieval phases.</p>
        <a href="/swagger" class="swagger-link">📖 View API Documentation (Swagger UI)</a>
    </div>

    <div class="card">
        <h2>Available Endpoints:</h2>
        <ul>
            <li><strong>GET</strong> <a href="/phase-one">/phase-one</a> - Get a single cat fact</li>
            <li><strong>GET</strong> <a href="/phase-two">/phase-two</a> - Get 5 cat facts (sequential)</li>
            <li><strong>GET</strong> <a href="/phase-three">/phase-three</a> - Get 10 cat facts (concurrent)</li>
            <li><strong>GET</strong> <a href="/headers">/headers</a> - Debug: View request headers</li>
        </ul>
    </div>

    <div class="card">
        <h2>Quick Links:</h2>
        <ul>
            <li><a href="/swagger">Interactive API Documentation (Swagger UI)</a></li>
            <li><a href="/swagger/spec">OpenAPI Specification (YAML)</a></li>
        </ul>
    </div>
</body>
</html>
`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

// Admin API handlers
func readyHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func aliveHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Alive")
}

// Start Admin API server
func startAdminServer() {
	adminMux := http.NewServeMux()
	adminMux.HandleFunc("/_/ready", readyHandler)
	adminMux.HandleFunc("/_/alive", aliveHandler)

	fmt.Println("🔧 Admin API available at http://localhost:11666")
	fmt.Println("   - Health check: http://localhost:11666/_/alive")
	fmt.Println("   - Readiness probe: http://localhost:11666/_/ready")

	if err := http.ListenAndServe(":11666", adminMux); err != nil {
		log.Printf("Error starting admin server: %v\n", err)
	}
}

func main() {
	// Start Admin API in a separate goroutine
	go startAdminServer()

	// API endpoints
	http.HandleFunc("/phase-one", phaseOneAPI)
	http.HandleFunc("/phase-two", phaseTwoAPI)
	http.HandleFunc("/phase-three", phaseThreeAPI)
	http.HandleFunc("/headers", headers)
	http.HandleFunc("/cat-facts", phaseFourAPI)

	// Documentation endpoints
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/swagger", swaggerUIHandler)
	http.HandleFunc("/swagger/", swaggerUIHandler) // Handle with trailing slash
	http.HandleFunc("/swagger/spec", swaggerSpecHandler)

	fmt.Println("🏠 Home page at http://localhost:8090")
	fmt.Println("📖 Swagger UI available at http://localhost:8090/swagger")

	if err := http.ListenAndServe(":8090", nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
