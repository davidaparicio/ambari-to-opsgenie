package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func check200(w http.ResponseWriter, r *http.Request) {
	log.Printf("200 req %s %s\n", r.Host, r.URL.Path)
	jsonFile, err := os.Open("api/examples/200.json")
	defer func() {
		if err := jsonFile.Close(); err != nil {
			log.Println(err)
		}
	}()
	if err != nil {
		log.Println(err)
	}
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(byteValue)
}

func check200critical(w http.ResponseWriter, r *http.Request) {
	log.Printf("200 critical called")
	jsonFile, err := os.Open("api/examples/200_critical.json")
	defer func() {
		if err := jsonFile.Close(); err != nil {
			log.Println(err)
		}
	}()
	if err != nil {
		log.Println(err)
	}
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(byteValue)
}

func check200warning(w http.ResponseWriter, r *http.Request) {
	log.Printf("200 warning called")
	jsonFile, err := os.Open("api/examples/200_warning.json")
	defer func() {
		if err := jsonFile.Close(); err != nil {
			log.Println(err)
		}
	}()
	if err != nil {
		log.Println(err)
	}
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(byteValue)
}

func check403(w http.ResponseWriter, r *http.Request) {
	log.Println("403 called")
	jsonFile, err := os.Open("api/examples/403.json")
	defer func() {
		if err := jsonFile.Close(); err != nil {
			log.Println(err)
		}
	}()
	if err != nil {
		log.Println(err)
	}
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	w.Write(byteValue)
}

func check500(w http.ResponseWriter, r *http.Request) {
	log.Println("500 called")
	jsonFile, err := os.Open("api/examples/500.json")
	defer func() {
		if err := jsonFile.Close(); err != nil {
			log.Println(err)
		}
	}()
	if err != nil {
		log.Println(err)
	}
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(byteValue)
}

func main() {
	r := http.NewServeMux()

	r.HandleFunc("/403", check403)
	r.HandleFunc("/500", check500)
	r.HandleFunc("/200", check200)
	r.HandleFunc("/200/critical", check200critical)
	r.HandleFunc("/200/warning", check200warning)

	srv := &http.Server{
		Addr:              ":1337",
		ReadTimeout:       1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		Handler:           r,
		//TLSConfig: tlsConfig,
	}

	log.Println("Server running on port 1337")
	if err := srv.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Printf("Server stopping...")
		} else {
			log.Fatal(err)
		}
	}
}
