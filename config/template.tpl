user nginx;
worker_processes 1;
pid /var/run/nginx.pid;

events {
    worker_connections 1024;
}

http {
    include       mime.types;
    default_type  application/octet-stream;

    sendfile             on;
    tcp_nopush           on;
    tcp_nodelay          on;
    keepalive_timeout    65;
    server_tokens        off;
    access_log           off;
    error_log            /dev/stderr;
    client_max_body_size 32M;

    server {
        listen       80;
        server_name  _;
    }

    {{ range $id, $container := .Containers }}
    server {
        listen       80;
        server_name  {{ $container.Name }}.{{ $.Domain }};

        location / {
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_pass http://{{ $container.InternalIp }};
            # proxy_pass http://{{ $container.ExternalIp }}:{{ $container.ExternalPort }};
        }
    }
    {{ end }}
}
