package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"
)

type TemplateHandler struct {
	dstFile          string
	templateFile     string
	containerHandler *ContainerHandler
}

var funcMap = template.FuncMap{
	"strip": func(s1, s2 string) string {
		return strings.Replace(s2, s1, "", -1)
	},
}

func (t *TemplateHandler) GenerateFile() {
	dstFile, err := os.Create(t.dstFile)
	defer dstFile.Close()

	if err != nil {
		log.Fatalln(err)
	}

	tplFileContents, err := ioutil.ReadFile(t.templateFile)
	if err != nil {
		log.Fatalln(err)
	}

	tpl := template.Must(template.New("vhosts").Funcs(funcMap).Parse(string(tplFileContents)))
	if err = tpl.Execute(dstFile, t.containerHandler.Proxy); err != nil {
		log.Fatalln(err)
	}

	log.Printf("Generated file '%s' from template '%s'", t.dstFile, t.templateFile)
}

func NewTemplateHandler(dstFile, templateFile string, containerHandler *ContainerHandler) *TemplateHandler {
	return &TemplateHandler{
		dstFile:          dstFile,
		templateFile:     templateFile,
		containerHandler: containerHandler,
	}
}
