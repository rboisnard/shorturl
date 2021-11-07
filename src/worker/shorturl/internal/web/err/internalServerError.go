package weberr

import (
	"bytes"
	"log"
	"net/http"
)

func InternalServerError(w http.ResponseWriter, r *http.Request, reason string) {
	var request bytes.Buffer
	r.Write(&request)
	log.Printf("error '%s' with '%s'", reason, request.String())
	http.Error(w, "unexpected behavior", http.StatusInternalServerError)
}
