package responder

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
)

// respClient contains the respCodes
type respClient struct {
	respCodes map[StatusCode]MessageDict
	ctx       *context.Context
}

// ResponseClient is the functions supported for responder
type ResponseClient interface {
	Error(w http.ResponseWriter, r *http.Request, err error)
	JSON(w http.ResponseWriter, r *http.Request, code int, data interface{})
}

// NewResponder will create a response client given the resp status dictionaries or map
func NewResponder(ctx *context.Context, respCodes map[StatusCode]MessageDict) ResponseClient {
	return &respClient{
		ctx:       ctx,
		respCodes: respCodes,
	}
}

// Respond is action on writing out response
func (e *respClient) respond(w http.ResponseWriter, httpStatusCode int, src interface{}) error {
	var body []byte
	var err error

	switch s := src.(type) {
	case []byte:
		if !json.Valid(s) {
			e.errorL(w, http.StatusInternalServerError, errors.New("Invalid JSON"))
			return err
		}
		body = s
	case string:
		body = []byte(s)
	default:
		if body, err = json.Marshal(src); err != nil {
			e.errorL(w, http.StatusInternalServerError, fmt.Errorf("failed to parse Json: %s", err.Error()))
			return err
		}
	}
	w.WriteHeader(httpStatusCode)
	_, err = w.Write(body)

	return err
}

// errorL will write given error to writer
func (e *respClient) errorL(w http.ResponseWriter, httpStatusCode int, err error) error {
	var body []byte
	if err != nil {
		body = []byte(err.Error())
	}

	w.WriteHeader(httpStatusCode)
	_, err = w.Write(body)

	return err
}

// Error is a wrapper for error for type ErrorMessage
func (e *respClient) Error(w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	em, ok := err.(ErrorMessage)
	if !ok {
		e.errorL(w, http.StatusInternalServerError, err)
		return
	}

	e.respond(w, em.MessageDict.Status, em)
}

// JSON is Respond wrapper for successful response
func (e *respClient) JSON(w http.ResponseWriter, r *http.Request, code int, data interface{}) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	e.respond(w, code, data)
}
