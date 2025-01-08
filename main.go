package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func echoHandler(w http.ResponseWriter, r *http.Request) {
	acceptEncodings := r.Header.Values("Accept-Encoding")
	if len(acceptEncodings) > 0 && acceptEncodings[0] == "gzip" {
		w.Header().Set("Content-Encoding", "gzip")
	}
	message := strings.TrimPrefix(r.URL.Path, "/echo/")
	if message == "" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(message))
}

func userAgentHandler(w http.ResponseWriter, r *http.Request) {
	userAgent := r.Header.Get("User-Agent")
	if userAgent == "" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(userAgent))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func fileHandler(w http.ResponseWriter, r *http.Request) {
	filePath := "tmp/" + strings.TrimPrefix(r.URL.Path, "/file/")
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()

		if err != nil {
			http.Error(w, "cant read request body", http.StatusBadRequest)
			return
		}

		err = os.WriteFile(filePath, body, 0644)
		if err != nil {
			http.Error(w, "cant write this file", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		return
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
	w.WriteHeader(http.StatusOK)
	w.Write(content)
}

func main() {
	http.HandleFunc("/user-agent", userAgentHandler)
	http.HandleFunc("/echo/", echoHandler)
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/file/", fileHandler)

	fmt.Println("start with port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server: ", err.Error())
		os.Exit(1)
	}
}
