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
			mockHeaders := make(http.Header)

			for _, header := range headers {
				mockHeaders.Add(header.key, header.value)
			}

			err = utils.WriteJSON(mrw, req, "testEnvelope", tt.data, 200, mockHeaders)
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

			h := mrw.Header().Get("Content-Type")

			if len(h) == 0 {
				t.Errorf("expected a value for header but got nil")
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

func Test_IsRequired(t *testing.T) {
	tests := []struct {
		Name     string
		Tag      string
		Data     any
		Expected bool
	}{
		{
			Name: "has tag",
			Tag:  "name",
			Data: struct {
				name string
			}{
				name: "bill",
			},
			Expected: true,
		},
		{
			Name: "does not have tag",
			Tag:  "acct",
			Data: struct {
				name string
			}{
				name: "bill",
			},
			Expected: false,
		},
		{
			Name: "tag is embedded struct",
			Tag:  "Age",
			Data: struct {
				Name          string
				SomethingElse struct {
					Age int
				}
			}{
				Name:          "bill",
				SomethingElse: struct{ Age int }{Age: 9},
			},
			Expected: true,
		},
		{
			Name: "tag is pointer",
			Tag:  "Age",
			Data: struct {
				Age *uint
			}{
				Age: func() *uint {
					age := uint(42)
					return &age
				}(),
			},
			Expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			u := NewUtils(app)

			hasField := u.IsRequired(tt.Data, tt.Tag)

			if tt.Expected && !hasField {
				t.Errorf("Expected %s, but got %v", tt.Tag, hasField)
			}
		})
	}
}

func Test_GenerateRandomBytes(t *testing.T) {
	u := NewUtils(app)

	hash, err := u.GenerateRandomBytes(uint32(u.ArgonParams.SaltLength))
	if err != nil {
		t.Fatal(err)
	}

	t.Log(hash)

	if len(hash) == 0 {
		t.Error("expected a hash with length, got no hash")
	}
}

func Test_NoLengthGenerateRandomBytes(t *testing.T) {
	u := NewUtils(app)
	hash, err := u.GenerateRandomBytes(0)
	if err != nil {
		t.Fatal(err)
	}

	if len(hash) != 0 {
		t.Error("expected a hash with length, got no hash")
	}
}

func Test_ValidatePhoneNumber(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		phoneNum string
		expected bool
	}{
		{
			name:     "Valid US Phone Number",
			phoneNum: "123-456-7890",
			expected: true,
		},
		{
			name:     "Valid US Phone Number with Extension",
			phoneNum: "123-456-7890 x123",
			expected: true,
		},
		{
			name:     "Valid International Phone Number",
			phoneNum: "+1 123-456-7890",
			expected: true,
		},
		{
			name:     "Invalid Phone Number - Too Short",
			phoneNum: "123-456-789",
			expected: false,
		},
		{
			name:     "Invalid Phone Number - Too Long",
			phoneNum: "123-456-78901",
			expected: false,
		},
		{
			name:     "Invalid Phone Number - Missing Digits",
			phoneNum: "123-456-789A",
			expected: false,
		},
		{
			name:     "Invalid Phone Number - Special Characters",
			phoneNum: "123-456-7890!",
			expected: false,
		},
	}

	// Initialize Utils
	u := NewUtils(app)

	// Run each test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := u.ValidatePhoneNumber(tt.phoneNum)
			if result != tt.expected {
				t.Errorf("ValidatePhoneNumber(%s) = %v; expected %v", tt.phoneNum, result, tt.expected)
			}
		})
	}
}
