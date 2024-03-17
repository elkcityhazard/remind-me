package errors

type Errors struct {
	ErrorLog map[string][]string
}

//		Add looks ip a key on the Error Log, and appends the value
//	 If the key doesn't exist, it creates it automatically
func (e *Errors) Add(key, value string) {
	e.ErrorLog[key] = append(e.ErrorLog[key], value)
}

// Get looks up  the error by key, and returns a slice of string
func (e *Errors) Get(key string) []string {
	return e.ErrorLog[key]
}

// IsValid checks the length of ErrorLog and returns true if there are no errors, indicating that everything is valid, and there are no errors
func (e *Errors) IsValid() bool {
	return len(e.ErrorLog) == 0
}
