package render

import (
	"encoding/json"
	"log"
	"net/http"
)

type JSON struct {
	Data any
}

func (j JSON) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, "application/json; charset=utf-8")
}

func (j JSON) Render(w http.ResponseWriter) (err error) {
	j.WriteContentType(w)
	rsp, err := json.Marshal(j.Data)
	if err != nil {
		log.Println("json marshal error:", err)
		return err
	}
	_, err = w.Write(rsp)
	if err != nil {
		log.Println("json write error:", err)
		return err
	}
	return
}
