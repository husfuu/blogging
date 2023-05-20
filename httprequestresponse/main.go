package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func yahalloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Yahallo world")
}

func yahalloHTML(w http.ResponseWriter, r *http.Request) {
	str :=
		`<html> 
			<head><title>Yahallo HTML</title></head>
			<body><h1>Yahallo HTML</h1></body>
		</html>
		`

	w.Write([]byte(str))
}

type Post struct {
	User    string
	Threads []string
}

func yahalloJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	post := &Post{
		User:    "Happy cat",
		Threads: []string{"happy1", "happy2", "happy3"},
	}
	json, _ := json.Marshal(post)
	w.Write(json)
}

func main() {
	server := http.Server{
		Addr: "127.0.0.1:8080",
	}

	http.HandleFunc("/yahallo", yahalloHandler)
	http.HandleFunc("/html", yahalloHTML)
	http.HandleFunc("/json", yahalloJSON)

	server.ListenAndServe()
}
