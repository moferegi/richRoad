package external

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type RequestParams struct {
	URL     string
	Method  string
	Headers map[string]string
	Body    interface{}
}

func HttpRequest[T any](params RequestParams) (*T, error) {
	var req *http.Request
	var err error

	if params.Method == "" {
		params.Method = "GET"
	}

	if params.Method == "GET" {
		// Convert the body to url.Values type
		data := url.Values{}
		for key, value := range params.Body.(map[string]interface{}) {
			data.Set(key, value.(string))
		}

		// Add the encoded body to the URL
		params.URL = params.URL + "?" + data.Encode()

		req, err = http.NewRequest(params.Method, params.URL, nil)
	} else {
		bodyBytes, perr := json.Marshal(params.Body)
		if perr != nil {
			return nil, perr
		}
		req, err = http.NewRequest(params.Method, params.URL, bytes.NewBuffer(bodyBytes))
	}

	if err != nil {
		return nil, err
	}

	// Set default header
	req.Header.Set("Content-Type", "application/json")

	for key, value := range params.Headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result T

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
