package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type GetSubtasks struct {
	UserId string `json:"userId"`
	Number int    `json:"numberOfSubtasks"`
}

type Subtask struct {
	Id     string   `json:"id"`
	Labels []string `json:"labels"`
}

type Subtasks struct {
	Subtasks []Subtask `json:"subtasks"`
}

func main() {
	pause := time.Second * 2
	var wg sync.WaitGroup
	url := "http://192.168.49.2/api/subtasks"
	goroutines := 10

	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		g := g
		go func() {
			userId := fmt.Sprintf("user-%d", g)

			for {
				// Prepare the request body
				requestBody := GetSubtasks{
					UserId: userId,
					Number: 10,
				}

				// Convert request body to JSON
				requestBodyJSON, err := json.Marshal(requestBody)
				if err != nil {
					fmt.Println("Error encoding request body:", err)
					return
				}

				// Perform GET request with query parameter in URL
				req, err := http.NewRequest("GET", url, bytes.NewBuffer(requestBodyJSON))
				if err != nil {
					fmt.Println("Error creating GET request:", err)
					return
				}

				// Send the request
				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					fmt.Println("Error making GET request:", err)
					return
				}

				// Read response body
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					fmt.Println("Error reading response body:", err)
					return
				}

				if resp.StatusCode != 200 {
					fmt.Println(string(body))
				}

				// Parse response data
				var responseData Subtasks
				if err := json.Unmarshal(body, &responseData); err != nil {
					fmt.Println("Error parsing response body:", err)
					return
				}

				err = resp.Body.Close()
				if err != nil {
					fmt.Println(err)
				}

				// Print parsed data
				for _, entry := range responseData.Subtasks {
					jsonData := map[string]interface{}{
						"userId":        userId,
						"assignedLabel": entry.Labels[rand.Intn(len(entry.Labels))],
					}

					jsonValue, err := json.Marshal(jsonData)
					if err != nil {
						fmt.Println("Error encoding JSON:", err)
						continue
					}

					req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/%s", url, entry.Id), bytes.NewBuffer(jsonValue))
					if err != nil {
						fmt.Println("Error making POST request:", err)
						continue
					}

					// Send the request
					client := &http.Client{}
					resp, err := client.Do(req)
					if err != nil {
						fmt.Println("Error making PATCH request:", err)
						return
					}

					if resp.StatusCode != 200 {
						fmt.Println(string(body))
					} else {
						fmt.Printf("classified subtask, outcome = %d\n", resp.StatusCode)
					}
				}
				time.Sleep(pause)
			}
		}()
	}

	wg.Wait()
}
