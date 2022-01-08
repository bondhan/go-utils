package responder

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
)

type respClient struct {
	respCodes map[StatusCode]MessageDict
}

type ResponseClient interface {
	Error(w http.ResponseWriter, r *http.Request, err error)
	JSON(w http.ResponseWriter, r *http.Request, code int, data interface{})
}

func NewResponder(respCodes map[StatusCode]MessageDict) ResponseClient {
	return &respClient{
		respCodes: respCodes,
	}
}

// Respond is response write to ResponseWriter
func (e *respClient) respond(w http.ResponseWriter, httpStatusCode int, src interface{}) {
	var body []byte
	var err error

	switch s := src.(type) {
	case []byte:
		if !json.Valid(s) {
			e.errorL(w, http.StatusInternalServerError, errors.New("Invalid JSON"))
			return
		}
		body = s
	case string:
		body = []byte(s)
	default:
		if body, err = json.Marshal(src); err != nil {
			e.errorL(w, http.StatusInternalServerError, fmt.Errorf("Failed to parse Json: %s", err.Error()))
			return
		}
	}
	w.WriteHeader(httpStatusCode)
	w.Write(body)
}

func (e *respClient) errorL(w http.ResponseWriter, httpStatusCode StatusCode, err error) {

	var em ErrorMessage
	em.MessageDict = e.respCodes[httpStatusCode]
	em.Message = err.Error()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	e.respond(w, em.MessageDict.Status, em)
}

// Error always receive err type ErrorMessage, other than that might produce error
func (e *respClient) Error(w http.ResponseWriter, r *http.Request, err error) {
	em, _ := err.(ErrorMessage)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	e.respond(w, em.MessageDict.Status, em)
}

// JSON is wrapped Respond when success response
func (e *respClient) JSON(w http.ResponseWriter, r *http.Request, code int, data interface{}) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	e.respond(w, code, data)
}
