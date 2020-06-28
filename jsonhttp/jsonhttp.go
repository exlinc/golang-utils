package jsonhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

// APIResponse contains the attributes found in an API response
type APIResponse struct {
	Message string      `json:"message"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Debug   string      `json:"debug,omitempty"`
}

// CheckableRequest defines an interface for request payloads that can be checked with the jsonhttp checker. See JSONDecodeAndCatchForAPI for the usage
type CheckableRequest interface {
	Parameters() error
}

// JSONSuccess returns a successful APIResponse on the http response with the provided parameters
func JSONSuccess(w http.ResponseWriter, data interface{}, message string) {
	if message == "" {
		message = "ok"
	}
	resp := APIResponse{
		Message: message,
		Success: true,
		Data:    data,
	}
	JSONWriter(w, resp, http.StatusOK)
}

// JSONError returns an APIResponse on the http response with the provided parameters and status code
func JSONError(w http.ResponseWriter, data interface{}, message string, debug string, statusCode int) {
	if message == "" {
		message = "error"
	}
	resp := APIResponse{
		Message: message,
		Success: false,
		Data:    data,
		Debug:   debug,
	}
	JSONWriter(w, resp, statusCode)
}

// JSONInternalError returns an internal server error APIResponse on the http response with the provided parameters
func JSONInternalError(w http.ResponseWriter, message string, debug string) {
	if message == "" {
		message = "error"
	}
	resp := APIResponse{
		Message: message,
		Success: false,
		Data:    nil,
		Debug:   debug,
	}
	JSONWriter(w, resp, http.StatusInternalServerError)
}

// JSONInternalError returns a bad request error APIResponse on the http response with the provided parameters
func JSONBadRequestError(w http.ResponseWriter, message string, debug string) {
	if message == "" {
		message = "bad_request"
	}
	resp := APIResponse{
		Message: message,
		Success: false,
		Data:    nil,
		Debug:   debug,
	}
	JSONWriter(w, resp, http.StatusBadRequest)
}

// JSONNotFoundError returns a not found error APIResponse on the http response with the provided parameters
func JSONNotFoundError(w http.ResponseWriter, message string, debug string) {
	if message == "" {
		message = "not_found"
	}
	resp := APIResponse{
		Message: message,
		Success: false,
		Data:    nil,
		Debug:   debug,
	}
	JSONWriter(w, resp, http.StatusNotFound)
}

// JSONForbiddenError returns a forbidden error APIResponse on the http response with the provided parameters
func JSONForbiddenError(w http.ResponseWriter, message string, debug string) {
	if message == "" {
		message = "forbidden"
	}
	resp := APIResponse{
		Message: message,
		Success: false,
		Data:    nil,
		Debug:   debug,
	}
	JSONWriter(w, resp, http.StatusForbidden)
}

// JSONDetailed returns the provided APIResponse on the http response with the provided HTTP status code
func JSONDetailed(w http.ResponseWriter, resp APIResponse, statusCode int) {
	if statusCode == 0 {
		statusCode = http.StatusOK
	}
	JSONWriter(w, resp, statusCode)
}

// JSONWriter provides a wrapper function to marshal an interface{} type to JSON and then send the bytes back over an http.ResponseWriter
func JSONWriter(w http.ResponseWriter, payload interface{}, statusCode int) {
	//dj, err := json.MarshalIndent(payload, "", "  ")
	dj, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		fmt.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, "%s", dj)
}

// JSONDecodeAndCatchForAPI is the primary function for decoding checkable (and non-checkable) payloads into structs. If the struct passed into `outStruct` satisfied the `CheckableRequest` interface, the check will also be run after decoding the JSON
func JSONDecodeAndCatchForAPI(w http.ResponseWriter, r *http.Request, outStruct interface{}) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&outStruct)
	if err != nil {
		JSONBadRequestError(w, "Invalid JSON", "")
		return err
	}
	if !isCheckableRequest(outStruct) {
		return nil
	}
	method := reflect.ValueOf(outStruct).MethodByName("Parameters").Interface().(func() error)
	err = method()
	if err != nil {
		JSONBadRequestError(w, "", err.Error())
		return err
	}
	return nil
}

// JSONDecodeNoCatch decodes checkable (and non-checkable) payloads into structs. If the struct passed into `outStruct` satisfied the `CheckableRequest` interface, the check will also be run after decoding the JSON
func JSONDecodeNoCatch(r *http.Request, outStruct interface{}) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&outStruct)
	if err != nil {
		return err
	}
	if !isCheckableRequest(outStruct) {
		return nil
	}
	method := reflect.ValueOf(outStruct).MethodByName("Parameters").Interface().(func() error)
	err = method()
	if err != nil {
		return err
	}
	return nil
}

func isCheckableRequest(checkAgainst interface{}) bool {
	reader := reflect.TypeOf((*CheckableRequest)(nil)).Elem()
	return reflect.TypeOf(checkAgainst).Implements(reader)
}
