package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

type Header struct {
	key   string
	value string
}

var headers = []Header{
	{
		key:   "Cache-Control",
		value: "no-store",
	},
	{
		key:   "Content-Type",
		value: "application/json;charset=utf-8",
	},
}

type MockResponseWriter struct {
	http.ResponseWriter
	Buffer  *bytes.Buffer
	Headers http.Header
}

func NewMockResponseWriter() *MockResponseWriter {
	return &MockResponseWriter{
		Buffer:  new(bytes.Buffer),
		Headers: make(http.Header),
	}
}

func (mrw *MockResponseWriter) Write(p []byte) (int, error) {
	return mrw.Buffer.Write(p)
}

func (mrw *MockResponseWriter) WriteHeader(statusCode int) {
}

func (mrw *MockResponseWriter) Header() http.Header {
	// Implement if needed for your tests
	return mrw.Headers
}

func Test_WriteJSON(t *testing.T) {
	utils := NewUtils(app)

	mrw := NewMockResponseWriter()

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		testType string
		data     any
	}{
		{
			testType: "fail",
			data:     make(chan bool),
		},
		{
			testType: "pass",
			data:     struct{ status string }{status: "okay"},
		},
	}

	// Call the method under test

	for _, tt := range tests {
		t.Run(tt.testType, func(t *testing.T) {
			err = utils.WriteJSON(mrw, req, "testEnvelope", tt.data)
			if err != nil {
				if tt.testType == "pass" {
					t.Fatalf("WriteJSON failed: %v", err)
				}
			}

			var result map[string]interface{}

			err = json.Unmarshal(mrw.Buffer.Bytes(), &result)
			if err != nil {
				if tt.testType == "pass" {
					t.Fatalf("Expected an error, but got none.")
				}
			}
		})
	}
}

func Test_ErrorJSON(t *testing.T) {
	utils := NewUtils(app)

	mrw := NewMockResponseWriter()

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		testType   string
		data       any
		statusCode int
	}{
		{
			testType:   "pass",
			data:       make(map[string]interface{}),
			statusCode: 200,
		},
		{
			testType:   "fail",
			data:       make(map[string]interface{}),
			statusCode: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testType, func(t *testing.T) {
			if tt.testType == "pass" {

				type msg struct {
					msg string
				}

				mockMsg := msg{
					msg: "hello world",
				}

				data := make(map[string]interface{})

				data["data"] = mockMsg

				tt.data = data
			} else {

				type msg struct {
					msg chan bool
				}

				mockMsg := msg{
					msg: make(chan bool),
				}

				data := make(map[string]interface{})

				data["data"] = mockMsg

				tt.data = data
			}

			mockHeaders := make(http.Header)

			for _, header := range headers {
				mockHeaders.Add(header.key, header.value)
			}

			err := utils.ErrorJSON(mrw, req, "error", tt.data, 400, mockHeaders)
			if err != nil {
				if tt.testType == "pass" {
					t.Errorf("expecting an error but got %v", err)
					return
				}
			}

			mrwHeaders := mrw.Header().Get("Cache-Control")

			if len(mrwHeaders) == 0 {
				t.Errorf("Expecting %s, but received no Content-Type header", "application/json;charset=utf-8")
			}

		})
	}
}
