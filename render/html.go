package render

import (
	"html/template"
	"log"
	"net/http"
)

type HTMLRender struct {
	Template *template.Template
}

type HTML struct {
	Data       any
	Name       string
	Template   *template.Template
	IsTemplate bool
}

func (h HTML) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, "text/html; charset=utf-8")
}

func (h HTML) Render(w http.ResponseWriter, code int) (err error) {
	h.WriteContentType(w)
	w.WriteHeader(code)
	if h.IsTemplate {
		err = h.Template.ExecuteTemplate(w, h.Name, h.Data)
	} else {
		_, err = w.Write([]byte(h.Data.(string)))
	}
	if err != nil {
		log.Println(err)
	}
	return
}
