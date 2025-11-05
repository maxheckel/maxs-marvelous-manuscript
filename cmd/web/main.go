package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/maxheckel/maxs-marvelous-manuscript/internal/api"
	"github.com/maxheckel/maxs-marvelous-manuscript/internal/db"
	"github.com/rs/cors"
)

const (
	defaultDataDir = "./data"
	defaultPort    = "8080"
)

func main() {
	// Get configuration from environment
	dataDir := getEnv("DATA_DIR", defaultDataDir)
	port := getEnv("PORT", defaultPort)

	// Ensure data directory exists
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	// Initialize database
	database, err := db.New(db.Config{
		DataDir: dataDir,
		DBName:  "dnd_assistant.db",
	})
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	recordingRepo := db.NewRecordingRepository(database)

	// Create API
	apiHandler := api.NewAPI(recordingRepo, dataDir)

	// Set up router
	router := mux.NewRouter()

	// Register API routes
	apiHandler.RegisterRoutes(router)

	// Serve static files from web/frontend/dist
	frontendDir := "./web/frontend/dist"
	if _, err := os.Stat(frontendDir); err == nil {
		router.PathPrefix("/").Handler(http.FileServer(http.Dir(frontendDir)))
	} else {
		// Fallback: serve a simple index page
		router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintf(w, `
<!DOCTYPE html>
<html>
<head>
    <title>D&D Session Assistant</title>
</head>
<body>
    <h1>D&D Session Assistant</h1>
    <p>Frontend not built yet. API is available at <a href="/api/health">/api/health</a></p>
</body>
</html>
`)
		})
	}

	// Set up CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000", "http://localhost:5173"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	handler := c.Handler(router)

	// Start server
	addr := fmt.Sprintf(":%s", port)
	fmt.Printf("🚀 D&D Session Assistant API\n")
	fmt.Printf("========================\n\n")
	fmt.Printf("Server starting on http://localhost%s\n", addr)
	fmt.Printf("API available at http://localhost%s/api\n", addr)
	fmt.Printf("Data directory: %s\n\n", filepath.Abs(dataDir))

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
