package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"sync"
	"time"
)

type CreationResponse struct {
	Id string `json:"id"`
}

const testDuration = time.Minute * 5

func main() {
	startTime := time.Now()

	pause := time.Second * 6

	var wg sync.WaitGroup
	url := "http://192.168.49.2/api/tasks"

	goroutines := 15

	go func() {
		wg.Add(1)
		/*
			TEST PATTERN
			change behaviour every 30 seconds
			- 0/30 		=> 6 sec
			- 30/60 	=> 5 sec
			- 60/90 	=> 4 sec
			- 90/120	=> 3 sec
			- 120/150	=> 2 sec
			- 150/180	=> 2 sec
			- 180/210	=> 3 sec
			- 210/240	=> 4 sec
			- 240/270	=> 6 sec
			- 270/300	=> 6 sec
		*/

		pauses := []int{6, 5, 4, 3, 2, 2, 3, 4, 6, 6}

		for startTime.Add(testDuration).After(time.Now()) {
			elapsedSeconds := int(time.Now().Sub(startTime).Seconds())

			if elapsedSeconds > 300 {
				elapsedSeconds = 299
			}

			pause = time.Duration(pauses[elapsedSeconds/30]) * time.Second
		}

		wg.Done()
	}()

	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		g := g
		go func() {
			i := 1

			for startTime.Add(testDuration).After(time.Now()) {
				time.Sleep(pause)

				jsonData := map[string]interface{}{
					"name":   fmt.Sprintf("Classification Job %d-%d", g, i),
					"labels": []string{"label1", "label2", "label3"},
				}

				jsonValue, err := json.Marshal(jsonData)
				if err != nil {
					fmt.Println("Error encoding JSON:", err)
					continue
				}

				resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
				if err != nil {
					fmt.Println("Error making POST request:", err)
					continue
				}

				body, err := io.ReadAll(resp.Body)
				if err != nil {
					fmt.Println("Error reading response body:", err)
					continue
				}

				err = resp.Body.Close()
				if err != nil {
					fmt.Println(err)
				}

				var responseData CreationResponse
				if err := json.Unmarshal(body, &responseData); err != nil {
					fmt.Println("Error parsing response body:", err)
					continue
				}

				// Send images (file needs to be inserted manually)
				file, err := os.Open(fmt.Sprintf("/tmp/few-images-%d.zip", g))
				if err != nil {
					fmt.Println("Error opening zip file:", err)
					continue
				}

				bodyBuf := &bytes.Buffer{}
				bodyWriter := multipart.NewWriter(bodyBuf)

				fileWriter, err := bodyWriter.CreateFormFile("images", fmt.Sprintf("/tmp/few-images-%d.zip", g))
				if err != nil {
					fmt.Println("Error writing zip file to body:", err)
					continue
				}

				_, err = io.Copy(fileWriter, file)
				if err != nil {
					fmt.Println("Error copying zip file to body:", err)
					continue
				}

				contentType := bodyWriter.FormDataContentType()
				err = bodyWriter.Close()
				if err != nil {
					fmt.Println(err)
				}

				err = file.Close()
				if err != nil {
					fmt.Println(err)
				}

				// Make multipart form data POST request
				resp, err = http.Post(fmt.Sprintf("%s/%s/upload", url, responseData.Id), contentType, bodyBuf)
				if err != nil {
					fmt.Println("Error making multipart form data POST request:", err)
					continue
				}

				// Read multipart form data response body
				body, err = io.ReadAll(resp.Body)
				if err != nil {
					fmt.Println("Error reading multipart form data response body:", err)
					continue
				}

				err = resp.Body.Close()
				if err != nil {
					fmt.Println(err)
				}

				fmt.Printf("client-%d: new task id = %s, sending images = %d\n", g, responseData.Id, resp.StatusCode)

				i++
			}

			wg.Done()
		}()
	}

	wg.Wait()
}
