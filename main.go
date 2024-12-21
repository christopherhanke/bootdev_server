package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/christopherhanke/bootdev_server/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits  atomic.Int32
	databaseQueries *database.Queries
	enviroment      string
	secret          string
}

func main() {
	const filepathRoot = "."
	const port = "8080"
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("loading enviroment failed: %v\n", err)
	}

	apiCfg.enviroment = os.Getenv("PLATFORM")
	apiCfg.secret = os.Getenv("SECRET")
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("open database failed: %v\n", err)
	}
	apiCfg.databaseQueries = database.New(db)

	// defining all access points with handlers
	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	mux.HandleFunc("GET /api/healthz", hanlderHealthz)

	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChip)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handlerDeleteChirp)

	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("PUT /api/users", apiCfg.handlerUpdateUser)

	mux.HandleFunc("POST /api/login", apiCfg.handlerLoginUser)

	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)

	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)

	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.handlerPolkaWebhooks)

	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port %s\n", filepathRoot, port)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("server failed: %v\n", err)
	}

}
