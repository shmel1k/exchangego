package errs

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var InternalError []byte = []byte(`{"status":500,"error":"internal error"}`)

type Error struct {
	Status int         `json:"status"`
	Err    interface{} `json:"error"`
}

func (e Error) Error() string {
	return fmt.Sprintf("%v", e.Err)
}

func WriteInternal(w http.ResponseWriter) {
	w.Write(InternalError)
	w.WriteHeader(http.StatusInternalServerError)
}

func WriteError(w http.ResponseWriter, err error) {
	if v, ok := err.(Error); ok {
		data, err1 := json.Marshal(v)
		if err1 != nil {
			WriteInternal(w)
			return
		}
		w.WriteHeader(v.Status)
		w.Write(data)
		return
	}
	WriteInternal(w)
}
