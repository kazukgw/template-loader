package templateloader

import (
	"html/template"
	"io/ioutil"
	fp "path/filepath"
)

type TemplateLoader struct {
	TemplateDir string
	Common      string
	GlobPattern string
	PreCompile  bool
	template.FuncMap
	Templates        map[string]*template.Template
	TemplateBuilders map[string]*TemplateBuilder
}

type TemplateBuilder struct {
	Paths []string
	template.FuncMap
}

func (tb *TemplateBuilder) Build() *template.Template {
	t := template.New("")
	t = t.Funcs(tb.FuncMap)
	for _, p := range tb.Paths {
		tmpl, err := t.ParseGlob(p)
		if err != nil {
			panic(err.Error())
		}
		t = tmpl
	}
	return t
}

func (tl *TemplateLoader) Get(name string) *template.Template {
	if tl.PreCompile {
		return tl.Templates[name]
	}
	return tl.TemplateBuilders[name].Build()
}

func (tl *TemplateLoader) Load() {
	tl.Templates = map[string]*template.Template{}
	tl.TemplateBuilders = map[string]*TemplateBuilder{}
	dir := tl.TemplateDir
	infos, err := ioutil.ReadDir(tl.TemplateDir)
	if err != nil {
		panic(err.Error())
	}
	commonp := fp.Join(dir, tl.Common)
	for _, info := range infos {
		if info.IsDir() && info.Name() != tl.Common {
			tb := &TemplateBuilder{
				FuncMap: tl.FuncMap,
				Paths: []string{
					fp.Join(commonp, tl.GlobPattern),
					fp.Join(dir, info.Name(), tl.GlobPattern),
				},
			}
			if tl.PreCompile {
				tl.Templates[info.Name()] = tb.Build()
			} else {
				tl.TemplateBuilders[info.Name()] = tb
			}
		}
	}
}
