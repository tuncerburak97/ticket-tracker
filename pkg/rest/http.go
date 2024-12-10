package rest

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type Request struct {
	Method  string
	URL     string
	Body    interface{}
	Headers map[string]interface{}
}

type Client struct {
	client *http.Client
}

var (
	once     sync.Once
	instance *Client
)

func GetRestClient() *Client {
	once.Do(func() {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
				MaxVersion: tls.VersionTLS12,
			},
		}
		instance = &Client{
			client: &http.Client{
				Timeout:   10 * time.Second,
				Transport: tr,
			},
		}
	})
	return instance
}
func (c *Client) SendHttpRequest(request Request) ([]byte, error) {
	jsonBody, err := json.Marshal(request.Body)
	if err != nil {
		return nil, err
	}

	httpRequest, err := http.NewRequest(request.Method, request.URL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	for name, value := range request.Headers {
		strValue, ok := value.(string)
		if ok {
			httpRequest.Header.Add(name, strValue)
		} else {
			return nil, fmt.Errorf("header value for %v is not a string", name)
		}
	}

	response, err := c.client.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v", err)
		}
	}(response.Body)

	if response.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, fmt.Errorf("unexpected status code: %v", response.StatusCode)
		}
		return nil, fmt.Errorf("unexpected status code: %v body: %v", response.StatusCode, string(bodyBytes))
	}

	return io.ReadAll(response.Body)
}
