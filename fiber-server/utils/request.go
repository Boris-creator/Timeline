package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type RequestOptions struct {
	Headers map[string]string
}

func Fetch[TResponse any](url string, method string, requestBody *map[string]any, options RequestOptions) (*TResponse, error) {
	requestJSON, _ := json.Marshal(requestBody)

	req, err := http.NewRequest(
		method,
		url,
		bytes.NewBuffer(requestJSON),
	)
	if err != nil {
		return nil, err
	}

	for key, value := range options.Headers {
		req.Header.Set(key, value)
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	responseRaw, _ := io.ReadAll(response.Body)
	var responseData TResponse
	json.Unmarshal(responseRaw, &responseData)
	return &responseData, nil
}
