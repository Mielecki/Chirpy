package main

import (
	"fmt"
	"net/http"
)

func main(){
	const filepathRoot = "."
	const port = "8080"

	cfg := apiConfig{}
	cfg.fileserverHits.Store(0)
	serveMux := http.NewServeMux()
	serveMux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	serveMux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	serveMux.HandleFunc("POST /admin/reset", cfg.handlerReset)
	serveMux.HandleFunc("GET /api/healthz", handlerReadiness)
	serveMux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)

	server := http.Server{Handler: serveMux, Addr: ":" + port}
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err.Error())
	}
}