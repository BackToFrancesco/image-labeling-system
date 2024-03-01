package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"sync"
	"time"
)

type CreationResponse struct {
	Id string `json:"id"`
}

const testDuration = time.Minute * 10

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
			change behaviour every minute
			- 0/1 	=> 6 sec
			- 1/2 	=> 5 sec
			- 2/3 	=> 4 sec
			- 3/4	=> 3 sec
			- 4/5	=> 2 sec
			- 5/6	=> 2 sec
			- 6/7	=> 3 sec
			- 7/8	=> 4 sec
			- 8/9	=> 6 sec
			- 9/10	=> 6 sec
		*/

		pauses := []int{6, 5, 4, 3, 2, 2, 3, 4, 6, 6}

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
			f, err := os.OpenFile(fmt.Sprintf("/tmp/client-%d.log", g), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Fatal(err)
			}

			defer func(f *os.File) {
				err := f.Close()
				if err != nil {
					log.Println(err)
				}
			}(f)

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

				t := time.Now()

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

				_, err = f.WriteString(fmt.Sprintf("%f,%f\n", time.Now().Sub(t).Seconds(), time.Now().Sub(startTime).Seconds()))
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
