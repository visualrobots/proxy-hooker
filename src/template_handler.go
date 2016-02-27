package main

import (
	"log"
	"os"
	"text/template"
)

type TemplateHandler struct {
	dstFile          string
	templateFile     string
	containerHandler *ContainerHandler
}

func (t *TemplateHandler) GenerateFile() {
	file, err := os.Create(t.dstFile)
	defer file.Close()

	if err != nil {
		log.Fatalln(err)
	}

	tpl, err := template.ParseFiles(t.templateFile)
	if err != nil {
		log.Fatalln(err)
	}

	if err = tpl.Execute(file, t.containerHandler.Proxy); err != nil {
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
