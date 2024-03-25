package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	go Run(&http.Server{
		Addr: ":8080",
	})

	testCases := []struct {
		method     string
		url        string
		body       string
		statusCode int
	}{
		{
			method:     http.MethodPost,
			url:        "/game",
			body:       `{"startDate":"2023-01-02T15:04:05.000Z","type":"chess","playerCount":2,"roudCount":2}`,
			statusCode: 200,
		},
		{
			method:     http.MethodGet,
			url:        "/game/1",
			body:       ``,
			statusCode: 200,
		},
		{
			method:     http.MethodDelete,
			url:        "/game/1",
			body:       ``,
			statusCode: 200,
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%s %s %+v -> %d", testCase.method, testCase.url, testCase.body, testCase.statusCode), func(t *testing.T) {
			req, err := http.NewRequest(testCase.method, "http://localhost:8080"+testCase.url, strings.NewReader(testCase.body))
			if err != nil {
				t.Fatalf("http.NewRequest failed: %v", err)
			}

			res, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("http.DefaultClient.Do failed: %v", err)
			}

			if res.StatusCode != http.StatusOK {
				bytes, err := io.ReadAll(res.Body)
				if err != nil {
					t.Fatalf("io.ReadAll res.Body failed: %v", err)
				}

				t.Fatalf("http.DefaultClient.Do failed with status %s -> %+v", res.Status, string(bytes))
			}
		})
	}
}
