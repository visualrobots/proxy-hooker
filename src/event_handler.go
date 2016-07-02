package main

import (
	"github.com/fsouza/go-dockerclient"
	"log"
)

const (
	EVENT_START = "start"
	EVENT_STOP  = "stop"
	EVENT_DIE   = "die"
)

type EventHandler struct {
	client           *docker.Client
	templateHandler  *TemplateHandler
	processHandler   *ProcessHandler
	containerHandler *ContainerHandler
	excludeContainer string
}

func (e *EventHandler) Listen() {
	events := make(chan *docker.APIEvents)
	e.client.AddEventListener(events)

	for {
		select {
		case event := <-events:
			switch event.Status {
			case EVENT_START:
				go e.HandleStartEvent(event.ID)
			case EVENT_STOP:
				fallthrough
			case EVENT_DIE:
				go e.HandleStopEvent(event.ID)
			}
		}
	}
}

func (e *EventHandler) HandleStartEvent(id string) {
	log.Printf("Received 'start' event for container '%s'", id)

	container := e.containerHandler.GetContainerInfo(id)

	if container.ExternalPort != "" && container.Name != e.excludeContainer {
		e.containerHandler.AddContainer(container)
		e.templateHandler.GenerateFile()
		e.processHandler.Reload()
	} else {
		log.Printf("Event ignored, no port exposed for container '%s' or filtered out", id)
	}
}

func (e *EventHandler) HandleStopEvent(id string) {
	log.Printf("Received 'stop' event for container '%s'", id)

	if e.containerHandler.RemoveContainer(id) {
		e.templateHandler.GenerateFile()
		e.processHandler.Reload()
	}
}

func NewEventHandler(client *docker.Client, templateHandler *TemplateHandler, processHandler *ProcessHandler, containerHandler *ContainerHandler, excludeContainer *string) *EventHandler {
	return &EventHandler{client, templateHandler, processHandler, containerHandler, *excludeContainer}
}
