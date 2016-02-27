package main

type Container struct {
	Id           string
	Name         string
	InternalIp   string
	ExternalIp   string
	ExternalPort string
}

type Proxy struct {
	Domain     string
	Containers map[string]*Container
}
