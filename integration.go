package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// Step 1: Start the server with a short delay to trigger "error" status
	var wg sync.WaitGroup
	wg.Add(1)
	server := NewServer(8080, 2*time.Second) // 2 seconds pending, 1 second completed, then error
	go server.Start(8080, &wg)
	
	// Allow the server to initialize
	time.Sleep(1 * time.Second) // Start after 1 second to see pending first
	
	// Step 2: Use the client library to poll the server
	client := NewClient("http://localhost:8080", 15, 2, 500*time.Millisecond) // Poll every 500ms to see pending, completed, and error
	finalStatus := client.WaitForCompletion()
	fmt.Printf("Final status: %s\n", finalStatus)
	
	// Step 3: Stop the server (in Go, this would be more complex, so we leave it running for simplicity)
	// Typically, we'd add graceful shutdown logic here
	fmt.Println("Test completed.")
	wg.Wait()
}
