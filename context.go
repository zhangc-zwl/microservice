package microservice

import (
	"html/template"
	"log"
	"net/http"
	"net/url"

	"github.com/zhangc-zwl/microservice/render"
)

type Context struct {
	W     http.ResponseWriter
	R     *http.Request
	engin *Engine
}

func (c *Context) HTML(status int, html string) (err error) {
	return c.Render(status, render.HTML{
		Data:       html,
		IsTemplate: false,
	})
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
	return c.Render(http.StatusOK, render.HTML{
		Data:       data,
		Name:       name,
		Template:   c.engin.HTMLRender.Template,
		IsTemplate: true,
	})
}

func (c *Context) JSON(status int, data any) error {
	return c.Render(status, render.JSON{
		Data: data,
	})
}

func (c *Context) XML(status int, data any) error {
	return c.Render(status, render.XML{
		Data: data,
	})
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

func (c *Context) FileFromFS(filepath string, fs http.FileSystem) {
	defer func(old string) {
		c.R.URL.Path = old
	}(c.R.URL.Path)

	c.R.URL.Path = filepath

	http.FileServer(fs).ServeHTTP(c.W, c.R)
}

func (c *Context) Redirect(status int, url string) {
	c.Render(status, render.Redirect{
		Code:     status,
		Request:  c.R,
		Location: url,
	})
}

func (c *Context) String(status int, format string, values ...any) (err error) {
	return c.Render(status, render.String{
		Format: format,
		Data:   values,
	})
}

func (c *Context) Render(code int, r render.Render) error {
	err := r.Render(c.W)
	c.W.WriteHeader(code)
	return err
}
