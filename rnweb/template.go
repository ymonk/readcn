package main

import (
    "net/http"
    "html/template"
    "path/filepath"
    "sync"
    "bytes"
)

const (
    templateDir = "./template"
    partialDir = templateDir + "/partial"
)

type templateHandler interface {
    String() (string, error)
    ServeHTTP(w http.ResponseWriter, r *http.Request)
    Data(key string) (interface{}, bool)
    SetData(key string, value interface{}) templateHandler
    AddDataMap(data map[string]interface{}) templateHandler
    ClearData() templateHandler
}

// New create a templateHandler. If development is true, returns a dynamic template
// which provide more development flexibility. Otherwise returns a static template
// that maximize the speed
func NewTemplateHandler(filename string, data map[string]interface{}, development bool) templateHandler {
    if data == nil {
        data = make(map[string]interface{})
    }
    if development {
        return newDynamicTemplateHandler(filename, data)
    } else {
        return newStaticTemplateHandler(filename, data)
    }
}

///////////////////////////////////////////////////////////////////////
// dynamicTemplateHandler - Flexible but slower, used in development.
///////////////////////////////////////////////////////////////////////

type dynamicTemplateHandler struct {
    once     sync.Once
    funcs    template.FuncMap
    filename string
    data     map[string]interface{}
}

func newDynamicTemplateHandler(filename string, data map[string]interface{}) *dynamicTemplateHandler {
    t := &dynamicTemplateHandler{
        filename: filename,
        data:     data,
    }
    t.once.Do(func() {
        t.funcs = template.FuncMap{
            "include": func(part string) (template.HTML, error) {
                buf := bytes.NewBuffer(nil)
                pt := template.Must(template.ParseFiles(filepath.Join(partialDir, part)))
                err := pt.Execute(buf, t.data)
                return template.HTML(buf.String()), err
            },
        }
    })
    return t
}

func (t *dynamicTemplateHandler) String() (string, error) {
    buf := bytes.NewBuffer(nil)
    tpl := template.Must(template.New(t.filename).Funcs(t.funcs).ParseFiles(filepath.Join(templateDir, t.filename)))
    if err := tpl.ExecuteTemplate(buf, t.filename, t.data); err != nil {
        return "", err
    }
    return buf.String(), nil
}

func (t *dynamicTemplateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    tpl := template.Must(template.New(t.filename).
    Funcs(t.funcs).
    ParseFiles(filepath.Join(templateDir, t.filename)))
    if err := tpl.ExecuteTemplate(w, t.filename, t.data); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
    t.ClearData()
}

// Data get template data
func (t *dynamicTemplateHandler) Data(key string) (interface{}, bool) {
    v, found := t.data[key]
    return v, found
}

// SetData set template data
func (t *dynamicTemplateHandler) SetData(key string, value interface{}) templateHandler {
    t.data[key] = value
    return t
}

// AddDataMap add template data map
func (t *dynamicTemplateHandler) AddDataMap(m map[string]interface{}) templateHandler {
    if len(t.data) == 0 {
        t.data = m
    } else {
        for k, v := range m {
            t.data[k] = v
        }
    }
    return t
}

// ClearData delete template data
func (t *dynamicTemplateHandler) ClearData() templateHandler {
    t.data = map[string]interface{}{}
    return t
}



//////////////////////////////////////////////////////////
// staticTemplateHandler - Faster, used in production.
//////////////////////////////////////////////////////////

// staticTemplateHandler handles HTTP request by rendering a template
type staticTemplateHandler struct {
    once     sync.Once
    filename string
    templ    *template.Template
    data     map[string]interface{}
}

// newStaticStaticTemplateHandler create a new Template Handler
func newStaticTemplateHandler(filename string, data map[string]interface{}) *staticTemplateHandler {
    if data == nil {
        data = map[string]interface{}{}
    }

    t := &staticTemplateHandler{
        filename: filename,
        data:     data,
    }

    t.once.Do(func() {
        funcs := template.FuncMap{
            "include": func(part string) (template.HTML, error) {
                buf := bytes.NewBuffer(nil)
                pt := template.Must(template.ParseFiles(filepath.Join(partialDir, part)))
                err := pt.Execute(buf, t.data)
                return template.HTML(buf.String()), err
            },
        }
        t.templ = template.Must(template.New(t.filename).Funcs(funcs).
        ParseFiles(filepath.Join(templateDir, t.filename)))
    })
    return t
}


func (t *staticTemplateHandler) String() (string, error) {
    buf := bytes.NewBuffer(nil)
    if err := t.templ.Execute(buf, t.data); err != nil {
        return "", err
    }
    return buf.String(), nil
}


// ServeHTTP provide HTTP service for staticTemplateHandler
func (t *staticTemplateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if err := t.templ.Execute(w, t.data); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
    t.ClearData()
}

// Data get template data
func (t *staticTemplateHandler) Data(key string) (interface{}, bool) {
    v, found := t.data[key]
    return v, found
}

// SetData set template data
func (t *staticTemplateHandler) SetData(key string, value interface{}) templateHandler {
    t.data[key] = value
    return t
}

// AddDataMap add template data map
func (t *staticTemplateHandler) AddDataMap(m map[string]interface{}) templateHandler {
    if len(t.data) == 0 {
        t.data = m
    } else {
        for k, v := range m {
            t.data[k] = v
        }
    }
    return t
}

// ClearData delete template data
func (t *staticTemplateHandler) ClearData() templateHandler {
    t.data = map[string]interface{}{}
    return t
}



















