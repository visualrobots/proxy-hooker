FROM gliderlabs/alpine:latest
RUN apk-install nginx bash s6

COPY proxy-hooker /usr/bin/
COPY config/template.tpl /etc/nginx/
COPY config/services.d /etc/services.d

VOLUME /etc/nginx

EXPOSE 80
ENTRYPOINT ["/bin/s6-svscan", "/etc/services.d"]
CMD []
