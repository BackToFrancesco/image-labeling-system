package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
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

const testDuration = time.Minute * 10

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
			change behaviour every minute
			- 0/1 	=> 9 sec
			- 1/2 	=> 7 sec
			- 2/3 	=> 5 sec
			- 3/4	=> 4 sec
			- 4/5	=> 3 sec
			- 5/6	=> 3 sec
			- 6/7	=> 4 sec
			- 7/8	=> 5 sec
			- 8/9	=> 9 sec
			- 9/10	=> 9 sec
		*/
		pauses := []int{9, 7, 5, 4, 3, 3, 4, 5, 9, 9}

		for startTime.Add(testDuration).After(time.Now()) {
			elapsedSeconds := int(time.Now().Sub(startTime).Seconds())

			if elapsedSeconds > 10*60 {
				elapsedSeconds = 10*60 - 1
			}

			pause = time.Duration(pauses[elapsedSeconds/60]) * time.Second
			time.Sleep(time.Second * 5)
		}

		wg.Done()
	}()

	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		g := g
		go func() {

			f, err := os.OpenFile(fmt.Sprintf("/tmp/classifier-%d.log", g), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Fatal(err)
			}

			defer func(f *os.File) {
				err := f.Close()
				if err != nil {
					log.Println(err)
				}
			}(f)

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

				t := time.Now()

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

				_, err = f.WriteString(fmt.Sprintf("%f,%f\n", time.Now().Sub(t).Seconds(), time.Now().Sub(startTime).Seconds()))
				if err != nil {
					fmt.Println(err)
				}
			}

			wg.Done()
		}()
	}

	wg.Wait()
}
