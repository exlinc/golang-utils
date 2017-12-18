# Utils for JSON APIs over HTTP

A suite of standardized functions for building REST APIs over HTTP

# Usage

## Request Payloads

The library is flexible here. They are simply passed into the `JSONDecodeAndCatchForAPI` function for parsing. For 'checkable' payloads, that struct must fulfil the `CheckableRequest` interface.

## Response Payloads

Response payloads are always of the `APIResponse` type:

```Go
// APIResponse contains the attributes found in an API response
type APIResponse struct {
	Message string      `json:"message"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Debug   string      `json:"debug,omitempty"`
}
```

## Using the response payload utils

Use the `JSON*` functions to send your messages, data, and debug through helpers that will package up `APIResponse` objects, marshal them to JSON, and then send them on your HTTP writer
