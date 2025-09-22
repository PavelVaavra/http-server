package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

func main() {
	const filepathRoot = "."
	const port = "8080"

	apiCfg := apiConfig{}
	mux := http.NewServeMux()

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", serverStatus)
	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsPrint)
	mux.HandleFunc("POST /admin/reset", apiCfg.metricsReset)
	mux.HandleFunc("POST /api/validate_chirp", validateChirp)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Printf("Serving files from %v on port: %v\n", filepathRoot, port)
	server.ListenAndServe()
}

func validateChirp(w http.ResponseWriter, r *http.Request) {
	const validLength = 140

	type chirp struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	ch := chirp{}
	err := decoder.Decode(&ch)
	if err != nil {
		handle500(w, r)
		return
	}

	if len(ch.Body) > validLength {
		handle400(w, r)
		return
	}

	type return200 struct {
		Valid bool `json:"valid"`
	}
	respBody := return200{
		Valid: true,
	}
	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(dat)
}

func handle500(w http.ResponseWriter, _ *http.Request) {
	type return500 struct {
		Error string `json:"error"`
	}
	respBody := return500{
		Error: "Something went wrong",
	}
	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)
	w.Write(dat)
}

func handle400(w http.ResponseWriter, _ *http.Request) {
	type return400 struct {
		Error string `json:"error"`
	}
	respBody := return400{
		Error: "Chirp is too long",
	}
	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)
	w.Write(dat)
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
	html := fmt.Sprintf(`<html>
<body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileserverHits.Load())
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

func (cfg *apiConfig) metricsReset(_ http.ResponseWriter, _ *http.Request) {
	cfg.fileserverHits.Store(0)
}
