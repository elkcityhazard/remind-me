package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/elkcityhazard/remind-me/internal/config"
	"golang.org/x/crypto/argon2"
)

var app *config.AppConfig

type Utilser interface {
	WriteJSON(w http.ResponseWriter, r *http.Request, envelope string, data interface{}) error
	ErrorJSON(w http.ResponseWriter, r *http.Request, enveloper string, data interface{}, statusCode int, headers ...http.Header) error
}

type argonParams struct {
	memory      int
	iterations  int
	parallelism int
	saltLength  int
	keylength   int
}

type Utils struct {
	app         *config.AppConfig
	maxFileSize int
	argonParams
}

// NewUtils is a utility helper to take care of certain tasks
// It needs the app config passed into it so it can have access
// to app wide items
func NewUtils(a *config.AppConfig) *Utils {
	app = a

	return &Utils{
		app:         app,
		maxFileSize: 1024 * 1024 * 1024 * 30,
		argonParams: argonParams{
			memory:      64 * 1024,
			iterations:  3,
			parallelism: 2,
			saltLength:  16,
			keylength:   32,
		},
	}
}

// WriteJSON takes in a responseWriter, request, enveolor, and data and write json to the response writer.
// it can return a potential error
// to pass in headers use headers := make(http.Header)
func (u *Utils) WriteJSON(w http.ResponseWriter, r *http.Request, envelope string, data interface{}, statusCode int, headers ...http.Header) error {
	payload := make(map[string]interface{})

	payload[envelope] = data

	w.Header().Set("Content-Type", "application/json;encoding=utf-8;")

	if len(headers) > 0 {
		for _, header := range headers {
			for key, value := range header {
				for _, v := range value {
					w.Header().Add(key, v)
				}
			}
		}
	}

	w.WriteHeader(statusCode)

	return json.NewEncoder(w).Encode(payload)
}

//  ErrorJSON is provided a response writer to write to, a pointer to a request, an envelope string, and some data
//  and returns an error json output.  It can potentially return an error

func (u *Utils) ErrorJSON(w http.ResponseWriter, r *http.Request, envelope string, data interface{}, statusCode int, headers ...http.Header) error {
	payload := map[string]interface{}{}

	payload[envelope] = data

	w.Header().Set("Content-Type", "application/json;encoding:utf-8;")

	for _, header := range headers {
		for key, values := range header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
	}

	w.WriteHeader(statusCode)

	return json.NewEncoder(w).Encode(payload)
}

func (u *Utils) GenerateSalt() ([]byte, error) {
	salt := make([]byte, 16)

	_, err := rand.Read(salt)

	if err != nil {
		return nil, err
	}

	return salt, nil
}

func (u *Utils) GenerateArgon2Password(password string, salt []byte) (string, error) {

	hash := argon2.IDKey([]byte(password), salt, uint32(u.argonParams.iterations), uint32(u.argonParams.memory), uint8(u.argonParams.parallelism), uint32(u.argonParams.keylength))

	encodedHash := base64.StdEncoding.EncodeToString(hash)

	return encodedHash, nil

}

func (u *Utils) VerifyArgon2Password(password string, storedHash string, storedSalt string) bool {
	hashBytes, err := base64.StdEncoding.DecodeString(storedHash)

	if err != nil {
		return false
	}

	saltBytes, err := base64.StdEncoding.DecodeString(storedSalt)

	if err != nil {
		return false
	}

	hash := argon2.IDKey([]byte(password), saltBytes, uint32(u.argonParams.iterations), uint32(u.argonParams.memory), uint8(u.argonParams.parallelism), uint32(u.argonParams.keylength))

	return subtle.ConstantTimeCompare(hash, hashBytes) == 1

}

// IsRequred traverse an interface, and looks to see if a key is present and returns a bool

func (u *Utils) IsRequired(s interface{}, key string) bool {

	val := reflect.ValueOf(s)

	// handle pointer
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// loop through the fields and determine if the field exists

	for i := 0; i < val.NumField(); i++ {
		// get the field based on index
		field := val.Field(i)
		//	get the field type
		fieldType := val.Type().Field(i)

		// check if the field name is equal to the key and return true if it is non zero
		if fieldType.Name == key {
			return !field.IsZero()
		}

		// use recursion and check if the field kind is a struct, and perform the same operation

		if field.Kind() == reflect.Struct {
			if u.IsRequired(field.Interface(), key) {
				return true
			}
		}
	}

	return false

}
