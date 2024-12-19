package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Client struct {
	BaseURL       string
	MaxRetries    int
	BackoffFactor int
	PollInterval  time.Duration
}

func NewClient(baseURL string, maxRetries int, backoffFactor int, pollInterval time.Duration) *Client {
	return &Client{
		BaseURL:       baseURL,
		MaxRetries:    maxRetries,
		BackoffFactor: backoffFactor,
		PollInterval:  pollInterval,
	}
}

func (c *Client) GetStatus() (string, error) {
	resp, err := http.Get(fmt.Sprintf("%s/status", c.BaseURL))
	if err != nil {
		return "error", fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()
	
	body, _ := ioutil.ReadAll(resp.Body)
	var result map[string]string
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "error", fmt.Errorf("failed to parse response: %v", err)
	}
	return result["result"], nil
}

func (c *Client) WaitForCompletion() string {
	lastStatus := ""
	for attempt := 1; attempt <= c.MaxRetries; attempt++ {
		status, err := c.GetStatus()
		if err != nil {
			fmt.Printf("Attempt %d: Error - %s\n", attempt, err)
			return "error"
		}
		fmt.Printf("Attempt %d: Status - %s\n", attempt, status)
		if status != lastStatus {
			fmt.Printf("Status changed to %s\n", status)
			lastStatus = status
		}
		if status == "completed" {
			fmt.Println("Successfully completed the process.")
			return "completed"
		} else if status == "error" {
			fmt.Println("An error occurred. Stopping retries.")
			return "error"
		}
		time.Sleep(c.PollInterval)
	}
	fmt.Println("Maximum retries reached. Timing out.")
	return "timeout"
}