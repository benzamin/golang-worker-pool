package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type BigTask struct{}

func NewBigTask() BigTask {
	return BigTask{}
}

func (e *BigTask) Run(payload Payload) JobResult {
	name := payload.Params["productName"]

	url := "https://dummyjson.com/products/add"

	bodyBytes := []byte(`{"title":"` + name + `"}`)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return NewJobResult(nil, err)
	}

	// Set the appropriate headers for the request
	req.Header.Add("Content-Type", "application/json")

	// Send the request to dummyJson, skipping SSL verify if required
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: transport, Timeout: time.Second * 60}
	resp, err := client.Do(req)
	if err != nil {
		return NewJobResult(nil, err)

	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return NewJobResult(nil, err)

	}

	var objmap map[string]interface{}
	err = json.Unmarshal(respBody, &objmap)
	if err != nil {
		return NewJobResult(nil, err)
	}
	//fmt.Print(objmap)
	return NewJobResult(objmap, nil)

}
