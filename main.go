package main

import (
	"log"
	"net/http"
)

func home(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("Games over !"))

}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	log.Println("launch web server on http://127.0.0.1:4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
