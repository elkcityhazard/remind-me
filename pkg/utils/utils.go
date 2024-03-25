package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/elkcityhazard/remind-me/internal/config"
	"golang.org/x/crypto/argon2"
)

var app *config.AppConfig

type Utilser interface {
	WriteJSON(w http.ResponseWriter, r *http.Request, envelope string, data interface{}) error
	ErrorJSON(w http.ResponseWriter, r *http.Request, enveloper string, data interface{}, statusCode int, headers ...http.Header) error
	CreateArgonHash(string, []byte) string
	GenerateRandomBytes(uint32) ([]byte, error)
	VerifyArgonHash(string, []byte) bool
}

type ArgonParams struct {
	Memory      int
	Iterations  int
	Parallelism int
	SaltLength  int
	KeyLength   int
	SaltKey     []byte
}

type Utils struct {
	app         *config.AppConfig
	maxFileSize int
	ArgonParams
}

// NewUtils is a utility helper to take care of certain tasks
// It needs the app config passed into it so it can have access
// to app wide items
func NewUtils(a *config.AppConfig) *Utils {
	app = a

	u := &Utils{
		app:         app,
		maxFileSize: 1024 * 1024 * 1024 * 30,
		ArgonParams: ArgonParams{
			Memory:      64 * 1024,
			Iterations:  3,
			Parallelism: 2,
			SaltLength:  6,
			KeyLength:   32,
		},
	}

	saltKey, err := u.GenerateRandomBytes(uint32(u.ArgonParams.SaltLength))
	if err != nil {
		panic(err)
	}

	u.ArgonParams.SaltKey = saltKey

	return u
}

//  GenerateRandomBytes returns a []byte of some randomness based on a salt length provided
//  it can return a slice of byte or error

func (u *Utils) GenerateRandomBytes(saltLength uint32) ([]byte, error) {
	//  make an empty byte slice based on saltLength

	b := make([]byte, saltLength)
	// fill byte slice with rand.Read
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

//  CreateArgonHash receives a plaintext password and a salt
//  and returns an encodedHashString

func (u *Utils) CreateArgonHash(plainTextPW string) string {
	hash := argon2.IDKey([]byte(plainTextPW), u.ArgonParams.SaltKey, uint32(u.Iterations), uint32(u.Memory), uint8(u.Parallelism), uint32(u.KeyLength))

	b64Salt := base64.RawStdEncoding.EncodeToString(u.ArgonParams.SaltKey)

	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, u.ArgonParams.Memory, u.ArgonParams.Iterations, u.ArgonParams.Parallelism, b64Salt, b64Hash)

	fmt.Println("Encoded Hash: ", encodedHash)

	return encodedHash
}

//   VerifyArgonHash accepts a plaintext password and compares it to a stored hash
//   and returns a bool (true) if it matches or (false) if it does not match

func (u *Utils) VerifyArgonHash(plaintextPW string, previousHash []byte) bool {
	hash := argon2.IDKey([]byte(plaintextPW), u.ArgonParams.SaltKey, uint32(u.ArgonParams.Iterations), uint32(u.ArgonParams.Memory), uint8(u.ArgonParams.Parallelism), uint32(u.ArgonParams.KeyLength))

	return subtle.ConstantTimeCompare(hash, previousHash) == 1
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
