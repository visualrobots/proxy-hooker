Proxy Hooker
============

Proxy Hooker is doing only one thing: **listen for Docker events and generate automatically a reverse proxy config file**.

- When a container is started, a new virtual host is created and the reverse proxy is reload gracefully
- When a container is stopped, the virtual host is removed

You can run containers that expose the port 80 under multiple virtual hosts without having to configure anything manually:

- myapp1.mydomain.tld -> myapp1 8001:80
- myapp2.mydomain.tld -> myapp2 8002:80
- myapp3.mydomain.tld -> myapp3 8003:80

The docker image use Nginx but any reverse proxy software can be used.

Features
--------

- Listen for container "start" and "stop" events
- Generate the reverse proxy file from a template (Golang template syntax)
- No dependency (written in Golang)

Author
------

Frédéric Guillot

License
-------

MIT

Installation
------------

Pull the Docker image:

```bash
docker pull fguillot/proxy-hooker
```

Try with your local Docker Machine:

```bash
docker run -d --name proxy-hooker \
    -p 80:80 \
    -v /var/lib/boot2docker:/certs:ro \
    -e DOCKER_HOST=$DOCKER_HOST \
    fguillot/proxy-hooker:latest
```

Examples
--------

Start a new container with automatic port assignation:

```bash
$ docker run -d --name my-container -P kanboard/kanboard:latest
1a214bfe2d3c1750561e7193e48ec219e246e4b57afe4939203a5d464b601456
```

Proxy Hooker will handle the event:

```bash
2016/02/27 20:01:04 Receive 'start' event for container '1a214bfe2d3c1750561e7193e48ec219e246e4b57afe4939203a5d464b601456'
2016/02/27 20:01:04 Adding container '1a214bfe2d3c1750561e7193e48ec219e246e4b57afe4939203a5d464b601456'
2016/02/27 20:01:04 Generated file '/etc/nginx/nginx.conf' from template '/etc/nginx/template.tpl'
2016/02/27 20:01:04 Command 'nginx -s reload' executed
```

The Nginx config file contains a new virtual host:

```
server {
    listen       80;
    server_name  my-container.mydomain.tld;

    location / {
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_pass http://172.17.0.6;
    }
}
```

Stop the container:

```bash
$ docker stop my-container
my-container
```

Proxy Hooker receive the event:

```
2016/02/27 20:01:14 Receive 'stop' event for container '1a214bfe2d3c1750561e7193e48ec219e246e4b57afe4939203a5d464b601456'
2016/02/27 20:01:14 Removing container '1a214bfe2d3c1750561e7193e48ec219e246e4b57afe4939203a5d464b601456'
2016/02/27 20:01:14 Generated file '/etc/nginx/nginx.conf' from template '/etc/nginx/template.tpl'
2016/02/27 20:01:14 Command 'nginx -s reload' executed
```

The virtual host is not there anymore.

Usage
-----

```bash
Usage of /usr/bin/proxy-hooker:
  -ca string
        TLS CA (default "/certs/ca.pem")
  -cert string
        TLS certificate (default "/certs/server.pem")
  -config string
        Config file generated (default "/etc/nginx/nginx.conf")
  -domain string
        Virtual host domain (default "mydomain.tld")
  -endpoint string
        Docker Host endpoint (default "tcp://192.168.99.100:2376")
  -key string
        TLS Key (default "/certs/server-key.pem")
  -reload-command string
        Command to run to reload the reverse proxy (default "nginx -s reload")
  -template string
        Configuration template (default "/etc/nginx/template.tpl")
```
