package microservice

import (
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/zhangc-zwl/microservice/render"
)

type Context struct {
	W          http.ResponseWriter
	R          *http.Request
	engin      *Engine
	queryCache url.Values
}

func (c *Context) HTML(status int, html string) (err error) {
	return c.Render(render.HTML{
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
	return c.Render(render.HTML{
		Data:       data,
		Name:       name,
		Template:   c.engin.HTMLRender.Template,
		IsTemplate: true,
	})
}

func (c *Context) JSON(status int, data any) error {
	return c.Render(render.JSON{
		Data: data,
	})
}

func (c *Context) XML(status int, data any) error {
	return c.Render(render.XML{
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
	c.Render(render.Redirect{
		Code:     status,
		Request:  c.R,
		Location: url,
	})
}

func (c *Context) String(status int, format string, values ...any) (err error) {
	return c.Render(render.String{
		Format: format,
		Data:   values,
	})
}

func (c *Context) Render(r render.Render) error {
	err := r.Render(c.W)
	return err
}

func (c *Context) DefaultQuery(key, defaultValue string) string {
	array, ok := c.GetQueryArray(key)
	if !ok {
		return defaultValue
	}
	return array[0]
}

func (c *Context) GetQuery(key string) string {
	c.initQueryCache()
	return c.queryCache.Get(key)
}

func (c *Context) QueryArray(key string) (values []string) {
	c.initQueryCache()
	values, _ = c.queryCache[key]
	return
}

func (c *Context) GetQueryArray(key string) (values []string, ok bool) {
	c.initQueryCache()
	values, ok = c.queryCache[key]
	return
}

func (c *Context) initQueryCache() {
	if c.R != nil {
		c.queryCache = c.R.URL.Query()
	} else {
		c.queryCache = url.Values{}
	}
}

func (c *Context) QueryMap(key string) (dict map[string]string) {
	dict, _ = c.GetQueryMap(key)
	return
}

func (c *Context) GetQueryMap(key string) (map[string]string, bool) {
	c.initQueryCache()
	return c.get(c.queryCache, key)
}

func (c *Context) get(m map[string][]string, key string) (map[string]string, bool) {
	//user[id]=1&user[name]=张三
	dict := make(map[string]string)
	exist := false
	for k, value := range m {
		if i := strings.IndexByte(k, '['); i >= 1 && k[0:i] == key {
			if j := strings.IndexByte(k[i+1:], ']'); j >= 1 {
				exist = true
				dict[k[i+1:][:j]] = value[0]
			}
		}
	}
	return dict, exist
}
