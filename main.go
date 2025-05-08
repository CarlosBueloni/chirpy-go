package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/CarlosBueloni/chirpy-go/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
}

func main() {
	const port = "8080"
	const filepathRoot = "/app/"
	godotenv.Load()

	dbUrl := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("Error opening sql database connection: %s", err)
	}

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		dbQueries:      database.New(db),
	}

	mux := http.NewServeMux()
	mux.Handle(filepathRoot, apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET  /api/healthz", handlerReadiness)
	mux.HandleFunc("POST  /api/validate_chirp", handlerValidate)

	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("GET  /admin/metrics", apiCfg.handlerMetrics)

	server := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
