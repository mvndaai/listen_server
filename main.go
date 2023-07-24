package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

func main() {
	http.HandleFunc("/", handler)
	port := ":8888"
	if v := os.Getenv("PORT"); v != "" {
		port = ":" + v
	}
	fmt.Println("Starting service on port", port)
	http.ListenAndServe(port, nil)
}

type printable struct {
	Path    string      `json:"path"`
	Method  string      `json:"method"`
	Body    string      `json:"body,omitempty"`
	Headers http.Header `json:"headers,omitempty"`
	Query   url.Values  `json:"query,omitempty"`
}

func handler(w http.ResponseWriter, req *http.Request) {
	p := printable{
		Path:    req.URL.Path,
		Method:  req.Method,
		Headers: req.Header,
		Query:   req.URL.Query(),
	}

	body, err := io.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(body) > 0 {
		p.Body = string(body)
	}

	b, err := json.MarshalIndent(p, "", "\t")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(b))
}
