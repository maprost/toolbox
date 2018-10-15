package net

import (
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
)

// CreateHTMLTemplates collect all html files inside the directory and subdirectories and put them into a template.
// Also use the relative path as key.
func (r *Router) CreateTemplatesFromPath(rootPath string) *template.Template {
	content := make(map[string]string)

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			key := path[len(rootPath):]

			b, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			content[key] = string(b)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	return r.CreateTemplates(content)
}

// CreateTemplates
func (r *Router) CreateTemplates(content map[string]string) *template.Template {
	mainTmpl := template.New("main")

	for key, text := range content {
		tmpl, err := loadTemplate(key, text, r.Config)
		if err != nil {
			panic(err)
		}
		_, err = mainTmpl.AddParseTree(key, tmpl.Tree)
		if err != nil {
			panic(err)
		}
	}

	r.engine.SetHTMLTemplate(mainTmpl)
	return mainTmpl
}

func loadTemplate(key string, text string, cfg *Config) (*template.Template, error) {
	tmpl, err := template.New(key).Delims(cfg.StartDelimiter, cfg.EndDelimiter).Parse(text)
	return tmpl, err
}
