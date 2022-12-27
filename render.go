package gan

import (
	"io"
	"log"
	"path/filepath"
	"text/template"
)

const layout = "templates/common/layout.html"

type Engine struct {
	params    map[string]any
	templates map[string]*template.Template
}

func NewEngine() *Engine {
	engine := &Engine{
		params:    make(map[string]any),
		templates: make(map[string]*template.Template),
	}
	engine.ParseTemplates()
	return engine
}

func (e *Engine) ParseTemplates() {
	files, err := filepath.Glob("templates/*.html")
	if err != nil {
		log.Fatal("Can not find templates")
	}
	tpl := template.New("main")
	for _, filename := range files {
		templateName := filepath.Base(filename)
		t, _ := tpl.Clone()
		e.templates[templateName] = template.Must(t.ParseFiles(layout, filename))
	}
}

func (e *Engine) Set(name string, value any) *Engine {
	e.params[name] = value
	return e
}

func merge(a map[string]any, b map[string]any) (output map[string]any) {
	output = map[string]any{}
	for k, v := range a {
		output[k] = v
	}
	for k, v := range b {
		output[k] = v
	}
	return output
}

func (e *Engine) Render(r io.Writer, name string, data map[string]any) error {
	tpl, ok := e.templates[name+".html"]
	if !ok {
		log.Fatal("[Template] The template does not exists", name)
	}
	return tpl.ExecuteTemplate(r, "base", merge(e.params, data))
}
