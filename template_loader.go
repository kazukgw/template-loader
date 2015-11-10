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
	template.FuncMap
	Templates map[string]*template.Template
}

func (tl *TemplateLoader) Load() {
	tl.Templates = map[string]*template.Template{}
	dir := tl.TemplateDir
	infos, err := ioutil.ReadDir(tl.TemplateDir)
	if err != nil {
		panic(err.Error())
	}
	commonp := fp.Join(dir, tl.Common)
	for _, info := range infos {
		if info.IsDir() && info.Name() != tl.Common {
			t, err := template.New("").Funcs(tl.FuncMap).ParseGlob(fp.Join(commonp, tl.GlobPattern))
			if err != nil {
				panic(err.Error())
			}
			p := fp.Join(dir, info.Name())
			t, err = t.ParseGlob(fp.Join(p, tl.GlobPattern))
			if err != nil {
				panic(err.Error())
			}
			tl.Templates[info.Name()] = template.Must(t, nil)
		}
	}
}
