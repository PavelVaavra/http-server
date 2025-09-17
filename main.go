package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

func main() {
	const filepathRoot = "."
	const port = "8080"

	apiCfg := apiConfig{}
	mux := http.NewServeMux()

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /healthz", serverStatus)
	mux.HandleFunc("GET /metrics", apiCfg.metricsPrint)
	mux.HandleFunc("POST /reset", apiCfg.metricsReset)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Printf("Serving files from %v on port: %v\n", filepathRoot, port)
	server.ListenAndServe()
}

func serverStatus(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) metricsPrint(w http.ResponseWriter, _ *http.Request) {
	s := fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())
	w.Write([]byte(s))
}

func (cfg *apiConfig) metricsReset(_ http.ResponseWriter, _ *http.Request) {
	cfg.fileserverHits.Store(0)
}
