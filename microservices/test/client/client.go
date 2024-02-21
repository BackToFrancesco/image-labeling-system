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

func main() {
	pause := time.Second * 2

	var wg sync.WaitGroup
	url := "http://192.168.49.2/api/tasks"

	goroutines := 10

	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func() {
			i := 1

			for {
				jsonData := map[string]interface{}{
					"name":   fmt.Sprintf("Classification Job %d", i),
					"labels": []string{"lab1", "lab2", "lab3"},
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
				file, err := os.Open("/tmp/images.zip")
				if err != nil {
					fmt.Println("Error opening zip file:", err)
					continue
				}

				bodyBuf := &bytes.Buffer{}
				bodyWriter := multipart.NewWriter(bodyBuf)

				fileWriter, err := bodyWriter.CreateFormFile("images", "images.zip")
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

				fmt.Printf("new task id = %s, sending images = %d\n", responseData.Id, resp.StatusCode)

				i++
				time.Sleep(pause)
			}
		}()
	}

	wg.Wait()
}
