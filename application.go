package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func root(w http.ResponseWriter, req *http.Request) {

	var filename string
	if filename = req.URL.Path[1:]; filename == "" {
		filename = "index.html"
	}
	filename = "public/" + filename
	fmt.Printf("Filename is: [%s]\n", filename)
	dat, err := ioutil.ReadFile(filename)

	acceptHeader := req.Header.Get("Accept")

	if strings.Contains(acceptHeader, "text/css") {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
	}

	if err != nil {
		io.WriteString(w, err.Error())
	} else {
		w.Write(dat)
	}
}

func pong(w http.ResponseWriter, req *http.Request) {
	log.Print("pong handler")
	io.WriteString(w, "pong v02")
}

func main() {
	http.HandleFunc("/", root)
	http.HandleFunc("/ping", pong)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
