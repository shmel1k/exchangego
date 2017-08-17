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
	data, err1 := json.Marshal(err)
	if err1 != nil {
		WriteInternal(w)
		return
	}
	w.Write(data)
}
