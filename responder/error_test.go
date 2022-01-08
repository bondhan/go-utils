package responder

import (
	"fmt"
	"github.com/guregu/null"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

const (
	ERR10001 StatusCode = 10001
	ERR10002 StatusCode = 10002
	ERR10003 StatusCode = 10003
)

var ErrorCodesExample = map[StatusCode]MessageDict{
	ERR10001: {Status: http.StatusBadRequest, Type: "DATA_INVALID", Code: ERR10001, Message: "Invalid Data Request"},
	ERR10002: {Status: http.StatusInternalServerError, Type: "INTERNAL_SERVER_ERROR", Code: ERR10002, Message: "Internal error"},
	ERR10003: {Status: http.StatusNotFound, Type: "NOT_FOUND", Code: ERR10003, Message: "Data Not Found"},
}

func TestComposeErrMsg(t *testing.T) {
	type args struct {
		emd MessageDict
		err error
	}

	tests := []struct {
		name    string
		args    args
		want    ErrorMessage
		isError bool
	}{
		{
			"Bad Request",
			args{
				ErrorCodesExample[ERR10001],
				errors.WithStack(errors.New("fail converting to integer")),
			},
			ErrorMessage{
				ErrorCodesExample[ERR10001],
				nil,
				nil,
			},
			true,
		},
		{
			"Internal Server",
			args{
				ErrorCodesExample[ERR10002],
				errors.WithStack(errors.New("internal server")),
			},
			ErrorMessage{
				ErrorCodesExample[ERR10002],
				nil,
				nil,
			},
			true,
		},
		{
			"Error is nil",
			args{
				ErrorCodesExample[ERR10001],
				nil,
			},
			ErrorMessage{
				ErrorCodesExample[ERR10001],
				nil,
				nil,
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ComposeErrMsg(tt.args.emd, tt.args.err)

			assert.Equal(t, got.Code, tt.want.Code)
			assert.Equal(t, got.Status, tt.want.Status)
			assert.Equal(t, got.Type, tt.want.Type)
			assert.Equal(t, got.Message, tt.want.Message)

			if tt.isError {
				assert.NotNil(t, got.Trace)
				assert.NotNil(t, got.ErrorMessage)
			} else {
				assert.Nil(t, got.Trace)
				assert.Nil(t, got.ErrorMessage)
			}
		})
	}
}

func TestErrorMessage_Error(t *testing.T) {
	type fields struct {
		MessageDict  MessageDict
		Trace        *string
		ErrorMessage *string
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"err 10001",
			fields{
				ErrorCodesExample[ERR10001],
				nil,
				nil,
			},
			fmt.Sprintf("%+v: %+v", ErrorCodesExample[ERR10001].Message, null.StringFromPtr(nil).ValueOrZero()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := ErrorMessage{
				MessageDict:  tt.fields.MessageDict,
				Trace:        tt.fields.Trace,
				ErrorMessage: tt.fields.ErrorMessage,
			}
			assert.Equalf(t, tt.want, e.Error(), "Error()")
		})
	}
}
