package main

import (
	"github.com/fsouza/go-dockerclient"
	"log"
	"sync"
)

type ContainerHandler struct {
	Proxy            *Proxy
	client           *docker.Client
	lock             sync.RWMutex
	excludeContainer string
}

func (c *ContainerHandler) BuildContainerList() {
	opts := docker.ListContainersOptions{All: true, Filters: map[string][]string{"status": []string{"running"}}}
	containers, err := c.client.ListContainers(opts)
	if err != nil {
		log.Printf("Unable to fetch list of containers: %s", err.Error())
	}

	for _, runningContainer := range containers {
		c.FilterContainer(c.GetContainerInfo(runningContainer.ID))
	}
}

func (c *ContainerHandler) GetContainerInfo(id string) *Container {
	container, _ := c.client.InspectContainer(id)
	externalIp, externalPort := c.findExternalNetworkSettings(container)
	internapIp, internaPort := c.findInternalNetworkSettings(container)

	return &Container{
		Id:           container.ID,
		Name:         container.Name[1:],
		ExternalPort: externalPort,
		ExternalIp:   externalIp,
		InternalIp:   internapIp,
		InternalPort: internaPort,
	}
}

func (c *ContainerHandler) FilterContainer(container *Container) bool {
	if container.InternalPort == "80" && container.Name != c.excludeContainer {
		c.AddContainer(container)
		return true
	}

	log.Printf("Container '%s' is filtered out", container.Name)
	return false
}

func (c *ContainerHandler) AddContainer(container *Container) {
	log.Printf("Adding container '%s'", container.Name)

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

func (c *ContainerHandler) findExternalNetworkSettings(container *docker.Container) (ip, port string) {
	if val, ok := container.NetworkSettings.Ports["80/tcp"]; ok && len(val) > 0 {
		port = val[0].HostPort
		ip = val[0].HostIP
	}

	return ip, port
}

func (c *ContainerHandler) findInternalNetworkSettings(container *docker.Container) (ip, port string) {
	ip = container.NetworkSettings.IPAddress

	if ip == "" {
		for _, network := range container.NetworkSettings.Networks {
			if network.IPAddress != "" {
				ip = network.IPAddress
				break
			}
		}
	}

	for mapping, _ := range container.NetworkSettings.Ports {
		if mapping.Proto() == "tcp" {
			port = mapping.Port()
			break
		}
	}

	return ip, port
}

func NewContainerHandler(client *docker.Client, domain string, excludeContainer string) *ContainerHandler {
	c := &ContainerHandler{
		client:           client,
		excludeContainer: excludeContainer,
		Proxy: &Proxy{
			Domain:     domain,
			Containers: make(map[string]*Container),
		},
	}

	return c
}
