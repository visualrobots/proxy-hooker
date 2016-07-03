package main

import (
	"flag"
	"github.com/fsouza/go-dockerclient"
	"log"
	"os"
)

func getDefaultValue(name, fallback string) string {
	value := os.Getenv(name)

	if value == "" {
		return fallback
	}

	return value
}

func main() {
	var reloadCommand = flag.String("command", getDefaultValue("COMMAND", "nginx -s reload"), "Command to run to reload the reverse proxy")
	var configFile = flag.String("config", getDefaultValue("CONFIG", "/etc/nginx/conf.d/vhosts.conf"), "Config file generated")
	var templateFile = flag.String("template", getDefaultValue("TEMPLATE", "/etc/nginx/template.tpl"), "Configuration template")
	var domain = flag.String("domain", getDefaultValue("DOMAIN", "mydomain.tld"), "Virtual host domain")
	var endpoint = flag.String("socket", getDefaultValue("SOCKET", "unix:///var/run/docker.sock"), "Docker Unix socket")
	var excludeContainer = flag.String("exclude", getDefaultValue("EXCLUDE", "proxy-hooker"), "Exclude a container name")
	flag.Parse()

	var client, err = docker.NewClient(*endpoint)
	if err != nil {
		log.Fatalln(err)
	}

	if err := client.Ping(); err != nil {
		log.Fatalln(err)
	}

	containerHandler := NewContainerHandler(client, *domain, *excludeContainer)
	containerHandler.BuildContainerList()

	templateHandler := NewTemplateHandler(*configFile, *templateFile, containerHandler)
	templateHandler.GenerateFile()

	processHandler := NewProcessHandler(*reloadCommand)
	processHandler.Reload()

	eventHandler := NewEventHandler(client, templateHandler, processHandler, containerHandler)
	eventHandler.Listen()
}
