package microservice

import (
	"encoding/json"
	"encoding/xml"
	"html/template"
	"log"
	"net/http"
)

type Context struct {
	W     http.ResponseWriter
	R     *http.Request
	engin *Engine
}

func (c *Context) HTML(status int, html string) {
	c.W.WriteHeader(status)
	c.W.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err := c.W.Write([]byte(html))
	if err != nil {
		log.Println(err)
	}
}

func (c *Context) HTMLTemplate(name string, funcMap template.FuncMap, data any, fileName ...string) {
	t := template.New(name)
	t.Funcs(funcMap)
	t, err := t.ParseFiles(fileName...)
	if err != nil {
		log.Println(err)
		return
	}
	c.W.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = t.Execute(c.W, data)
	if err != nil {
		log.Println(err)
	}
}

func (c *Context) HTMLTemplateGlob(name string, funcMap template.FuncMap, pattern string, data any) {
	t := template.New(name)
	t.Funcs(funcMap)
	t, err := t.ParseGlob(pattern)
	if err != nil {
		log.Println(err)
		return
	}
	c.W.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = t.Execute(c.W, data)
	if err != nil {
		log.Println(err)
	}
}

func (c *Context) Template(name string, data any) error {
	c.W.Header().Set("Content-Type", "text/html; charset=utf-8")
	template := c.engin.HTMLRender.Template
	err := template.ExecuteTemplate(c.W, name, data)
	if err != nil {
		log.Println(err)
	}
	return err
}

func (c *Context) JSON(status int, data any) error {
	c.W.Header().Set("Content-Type", "application/json; charset=utf-8")
	c.W.WriteHeader(status)
	rsp, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = c.W.Write(rsp)
	if err != nil {
		return err
	}
	return nil
}

func (c *Context) XML(status int, data any) error {
	header := c.W.Header()
	header["Content-Type"] = []string{"application/xml; charset=utf-8"}
	c.W.WriteHeader(status)
	err := xml.NewEncoder(c.W).Encode(data)
	if err != nil {
		return err
	}
	return nil
}
