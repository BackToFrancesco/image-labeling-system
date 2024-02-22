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

const testDuration = time.Minute * 5

func main() {
	startTime := time.Now()

	pause := time.Second * 9
	var wg sync.WaitGroup
	url := "http://192.168.49.2/api/subtasks"

	goroutines := 10

	go func() {
		wg.Add(1)
		/*
			TEST PATTERN
			change behaviour every 30 seconds
			- 0/30 		=> 9 sec
			- 30/60 	=> 7 sec
			- 60/90 	=> 5 sec
			- 90/120	=> 4 sec
			- 120/150	=> 3 sec
			- 150/180	=> 3 sec
			- 180/210	=> 4 sec
			- 210/240	=> 5 sec
			- 240/270	=> 9 sec
			- 270/300	=> 9 sec
		*/
		pauses := []int{9, 7, 5, 4, 3, 3, 4, 5, 9, 9}

		for startTime.Add(testDuration).After(time.Now()) {
			elapsedSeconds := int(time.Now().Sub(startTime).Seconds())

			pause = time.Duration(pauses[elapsedSeconds/30]) * time.Second
			time.Sleep(time.Second * 3)

			if elapsedSeconds > 300 {
				elapsedSeconds = 299
			}
		}

		wg.Done()
	}()

	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		g := g
		go func() {
			userId := fmt.Sprintf("user-%d", g)

			for startTime.Add(testDuration).After(time.Now()) {
				time.Sleep(pause)

				// Prepare the request body
				requestBody := GetSubtasks{
					UserId: userId,
					Number: 5,
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
						fmt.Printf("classifier-%d: classified subtask, outcome = %d\n", g, resp.StatusCode)
					}
				}
			}

			wg.Done()
		}()
	}

	wg.Wait()
}
