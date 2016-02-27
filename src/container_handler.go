package main

import (
	"github.com/fsouza/go-dockerclient"
	"log"
	"sync"
)

type ContainerHandler struct {
	Proxy  *Proxy
	client *docker.Client
	lock   sync.RWMutex
}

func (c *ContainerHandler) GetContainerInfo(id string) *Container {
	container, _ := c.client.InspectContainer(id)
	port := ""
	ip := ""

	if val, ok := container.NetworkSettings.Ports["80/tcp"]; ok && len(val) > 0 {
		port = val[0].HostPort
		ip = val[0].HostIP
	}

	return &Container{
		Id:           container.ID,
		Name:         container.Name[1:],
		ExternalPort: port,
		ExternalIp:   ip,
		InternalIp:   container.NetworkSettings.IPAddress,
	}
}

func (c *ContainerHandler) AddContainer(container *Container) {
	log.Printf("Adding container '%s'", container.Id)

	c.lock.Lock()
	defer c.lock.Unlock()

	c.Proxy.Containers[container.Id] = container
}

func (c *ContainerHandler) RemoveContainer(id string) bool {
	if _, ok := c.Proxy.Containers[id]; ok {
		log.Printf("Removing container '%s'", id)

		c.lock.Lock()
		defer c.lock.Unlock()

		delete(c.Proxy.Containers, id)
		return true
	}

	return false
}

func NewContainerHandler(client *docker.Client, domain string) *ContainerHandler {
	c := &ContainerHandler{
		client: client,
		Proxy: &Proxy{
			Domain:     domain,
			Containers: make(map[string]*Container),
		},
	}

	return c
}
