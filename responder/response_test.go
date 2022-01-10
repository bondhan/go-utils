package responder

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var nr ResponseClient

func TestMain(m *testing.M) {
	nr = NewResponder(nil, ErrorCodesExample)
}

func TestNewResponder(t *testing.T) {
	assert.NotNil(t, nr)
	_, ok := nr.(ResponseClient)
	if !ok {
		t.Errorf("object created no equal")
	}
}
