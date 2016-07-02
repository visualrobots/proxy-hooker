package main

import (
	"flag"
	"github.com/fsouza/go-dockerclient"
	"log"
)

func main() {
	var reloadCommand = flag.String("reload-command", "nginx -s reload", "Command to run to reload the reverse proxy")
	var configFile = flag.String("config", "/etc/nginx/nginx.conf", "Config file generated")
	var templateFile = flag.String("template", "/etc/nginx/template.tpl", "Configuration template")
	var domain = flag.String("domain", "mydomain.tld", "Virtual host domain")
	var endpoint = flag.String("socket", "unix:///var/run/docker.sock", "Docker Unix socket")
	flag.Parse()

	var client, err = docker.NewClient(*endpoint)
	if err != nil {
		log.Fatalln(err)
	}

	if err := client.Ping(); err != nil {
		log.Fatalln(err)
	}

	containerHandler := NewContainerHandler(client, *domain)
	processHandler := NewProcessHandler(*reloadCommand)

	templateHandler := NewTemplateHandler(*configFile, *templateFile, containerHandler)
	templateHandler.GenerateFile()

	eventHandler := NewEventHandler(client, templateHandler, processHandler, containerHandler)
	eventHandler.Listen()
}
