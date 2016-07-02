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
