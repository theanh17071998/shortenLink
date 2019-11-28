package main

import (
	"controller"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	//Tạo router
	router := mux.NewRouter()
	fs := http.FileServer(http.Dir("static"))
	router.Handle("/", fs).Methods("GET")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	router.HandleFunc("/", controller.Shorten).Methods("POST")
	router.HandleFunc("/{shortcode}", controller.Redirect).Methods("GET")

	var port string
	if port = os.Getenv("PORT"); port == "" {
		port = "8686"
	}
	fmt.Println("Server is running on port:", port)
	srv := &http.Server{
		Handler:      router,
		Addr:         "0.0.0.0:" + port,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	srv.SetKeepAlivesEnabled(false)

	log.Fatal(srv.ListenAndServe())
}
