package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {

	addr := flag.String("addr", "8080", "HTTP network address")
	flag.Parse()

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	mux.HandleFunc("/", homePage)
	mux.HandleFunc("/ascii-art/", AsciiArt)
	mux.HandleFunc("/download/", ExportFile)

	log.Printf("Server is listening... http://localhost:%s", *addr)

	err := http.ListenAndServe(":"+*addr, mux)
	if err != nil {
		log.Fatal(err)
	}
}
