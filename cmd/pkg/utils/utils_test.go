package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

type MockResponseWriter struct {
	http.ResponseWriter
	Buffer *bytes.Buffer
}

func NewMockResponseWriter() *MockResponseWriter {
	return &MockResponseWriter{
		Buffer: new(bytes.Buffer),
	}
}

func (mrw *MockResponseWriter) Write(p []byte) (int, error) {
	return mrw.Buffer.Write(p)
}

func (mrw *MockResponseWriter) WriteHeader(statusCode int) {

}

func (mrw *MockResponseWriter) Header() http.Header {
	// Implement if needed for your tests
	return http.Header{}
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
