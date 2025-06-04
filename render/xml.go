package render

import (
	"encoding/xml"
	"log"
	"net/http"
)

type XML struct {
	Data any
}

func (x XML) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, "application/xml; charset=utf-8")
}

func (x XML) Render(w http.ResponseWriter, code int) (err error) {
	x.WriteContentType(w)
	w.WriteHeader(code)
	if err = xml.NewEncoder(w).Encode(x.Data); err != nil {
		log.Println("xml render error:", err)
		return err
	}
	return
}
