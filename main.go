package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

func main() {
	http.HandleFunc("/user-agent", userAgentHandler)
	http.HandleFunc("/echo/", echoHandler)
	http.HandleFunc("/", rootHandler)

	fmt.Println("start with port 6969")
	if err := http.ListenAndServe(":6969", nil); err != nil {
		fmt.Println("Error starting server: ", err.Error())
		os.Exit(1)
	}
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
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
