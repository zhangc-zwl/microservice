package microservice

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"unicode"
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

func (c *Context) HTMLTemplate(name string, data any, fileName ...string) (err error) {
	t := template.New(name)
	t, err = t.ParseFiles(fileName...)
	if err != nil {
		log.Println(err)
		return
	}
	c.W.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = t.Execute(c.W, data)
	if err != nil {
		log.Println(err)
	}
	return
}

func (c *Context) HTMLTemplateGlob(name string, pattern string, data any) (err error) {
	t := template.New(name)
	t, err = t.ParseGlob(pattern)
	if err != nil {
		log.Println(err)
		return
	}
	c.W.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = t.Execute(c.W, data)
	if err != nil {
		log.Println(err)
	}
	return
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

func (c *Context) File(fileName string) {
	http.ServeFile(c.W, c.R, fileName)
}

func (c *Context) FileAttachment(filepath, filename string) {
	if isASCII(filename) {
		c.W.Header().Set("Content-Disposition", `attachment; filename="`+filename+`"`)
	} else {
		c.W.Header().Set("Content-Disposition", `attachment; filename*=UTF-8''`+url.QueryEscape(filename))
	}
	http.ServeFile(c.W, c.R, filepath)
}

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}

func (c *Context) FileFromFS(filepath string, fs http.FileSystem) {
	defer func(old string) {
		c.R.URL.Path = old
	}(c.R.URL.Path)

	c.R.URL.Path = filepath

	http.FileServer(fs).ServeHTTP(c.W, c.R)
}

func (c *Context) Redirect(status int, url string) {
	if (status < http.StatusMultipleChoices || status > http.StatusPermanentRedirect) && status != http.StatusCreated {
		panic(fmt.Sprintf("Cannot redirect with status code %d", status))
	}
	http.Redirect(c.W, c.R, url, status)
}
