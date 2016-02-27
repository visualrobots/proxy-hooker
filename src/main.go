package main

import (
	"flag"
	"github.com/fsouza/go-dockerclient"
	"log"
	"os"
)

func main() {
	var reloadCommand = flag.String("reload-command", "nginx -s reload", "Command to run to reload the reverse proxy")
	var configFile = flag.String("config", "/etc/nginx/nginx.conf", "Config file generated")
	var templateFile = flag.String("template", "/etc/nginx/template.tpl", "Configuration template")
	var domain = flag.String("domain", "mydomain.tld", "Virtual host domain")
	var endpoint = flag.String("endpoint", os.Getenv("DOCKER_HOST"), "Docker Host endpoint")
	var cert = flag.String("cert", "/certs/server.pem", "TLS certificate")
	var key = flag.String("key", "/certs/server-key.pem", "TLS Key")
	var ca = flag.String("ca", "/certs/ca.pem", "TLS CA")
	flag.Parse()

	var client, err = docker.NewTLSClient(*endpoint, *cert, *key, *ca)
	if err != nil {
		log.Fatalln(err)
	}

	if err := client.Ping(); err != nil {
		log.Fatalln(err)
	}

	containerHandler := NewContainerHandler(client, *domain)
	processHandler := NewProcessHandler(*reloadCommand)
	templateHandler := NewTemplateHandler(*configFile, *templateFile, containerHandler)
	eventHandler := NewEventHandler(client, templateHandler, processHandler, containerHandler)
	eventHandler.Listen()
}
