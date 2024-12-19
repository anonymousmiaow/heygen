package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	StartTime time.Time
	DelayTime time.Duration
	mux       *http.ServeMux
}

func NewServer(port int, delayTime time.Duration) *Server {
	server := &Server{
		StartTime: time.Now(),
		DelayTime: delayTime,
		mux:       http.NewServeMux(),
	}
	server.mux.HandleFunc("/status", server.handleStatus)
	return server
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	currentTime := time.Since(s.StartTime)
	var result string
	if currentTime < s.DelayTime {
		result = "pending"
	} else if currentTime < s.DelayTime+1*time.Second { // Completed state lasts 1 second instead of 500ms
		result = "completed"
	} else {
		result = "error"
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"result": result})
}

func (s *Server) Start(port int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Starting server on port %d with delay of %v seconds.\n", port, s.DelayTime.Seconds())
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), s.mux)
	if err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}