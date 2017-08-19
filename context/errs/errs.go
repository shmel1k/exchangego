package errs

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var InternalError []byte = []byte(`{"status":500,"error":"internal error"}`)

var ErrUserNotExists = Error{
	Status: http.StatusForbidden,
	Err:    "no user",
}

type Error struct {
	Status int         `json:"status"`
	Err    interface{} `json:"error"`
}

func (e Error) String() string {
	return fmt.Sprintf("%v", e.Err)
}

func (e Error) Error() string {
	return fmt.Sprintf("%v", e.Err)
}

func WriteInternal(w http.ResponseWriter) {
	w.Write(InternalError)
	w.WriteHeader(http.StatusInternalServerError)
}

func WriteError(w http.ResponseWriter, err error) int {
	if v, ok := err.(Error); ok {
		data, err1 := json.Marshal(v)
		if err1 != nil {
			WriteInternal(w)
			return http.StatusInternalServerError
		}
		w.WriteHeader(v.Status)
		w.Write(data)
		return v.Status
	}
	WriteInternal(w)
	return http.StatusInternalServerError
}
