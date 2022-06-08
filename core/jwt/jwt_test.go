package jwt

import (
	"errors"
	"testing"
)

func Test_registerJWT_ExpiredTokenError(t *testing.T) {
	j := NewJWT("123")
	err1 := j.ExpiredTokenError()
	err2 := j.ExpiredTokenError()
	if errors.As(err1, &err2) {
		t.Log("err1 == err2")
	}

	return
}
