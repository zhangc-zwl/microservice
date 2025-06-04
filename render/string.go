package render

import (
	"fmt"
	"log"
	"net/http"

	"github.com/zhangc-zwl/microservice/internal/bytesconv"
)

type String struct {
	Format string
	Data   []any
}

func (s String) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, "text/plain; charset=utf-8")
}

func (s String) Render(w http.ResponseWriter) (err error) {
	s.WriteContentType(w)
	if len(s.Data) > 0 {
		_, err = fmt.Fprintf(w, s.Format, s.Data...)
		return
	}
	_, err = w.Write(bytesconv.StringToBytes(s.Format))
	if err != nil {
		log.Println("render string error:", err)
	}
	return
}
