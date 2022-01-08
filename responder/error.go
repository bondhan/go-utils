// Package responder is generic to send response when success or error
// The format when success is:
//		{
//			"meta": {
//				"page":	1,
//				"perPage":	10,
//				"totalPage": 100
//			},
//			"data": {
//				"title": "a title",
//				...
//			}
//		}
//
// The format when error is:
//		{
//			"code":	12345,
//			"type":	"INVALID REQUEST",
//			"message": "value must be numeric",
//			"errMessage": "cannot convert err (type error) to type flags.Error"
//		}
//
//
//

package responder

import (
	"encoding/json"
	"fmt"

	"github.com/guregu/null"
)

type StatusCode int

type Message struct {
	Meta interface{} `json:"meta,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

func (i *Message) MarshalBinary() ([]byte, error) {
	return json.Marshal(i)
}

type MessageDict struct {
	Status  int        `json:"-"`
	Type    string     `json:"type"`
	Code    StatusCode `json:"code"`
	Message string     `json:"message"`
}

type ErrorMessage struct {
	MessageDict
	Trace        *string `json:"-"`
	ErrorMessage *string `json:"errorMessage,omitempty"`
}

func (e ErrorMessage) Error() string {
	return fmt.Sprintf("%+v: %+v", e.Message, null.StringFromPtr(e.Trace).ValueOrZero())
}

// ComposeErrMsg ..
func ComposeErrMsg(emd MessageDict, err error) ErrorMessage {

	var errMsg *string
	var errTrace *string

	if err != nil {
		s := err.Error()
		errMsg = &s

		ss := fmt.Sprintf("%+v", err)
		errTrace = &ss
	}

	em := ErrorMessage{
		emd,
		errTrace,
		errMsg,
	}

	return em
}
