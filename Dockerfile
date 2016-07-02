FROM alpine:3.4
RUN apk update && apk add nginx bash s6 curl && rm -rf /var/cache/apk/*

COPY proxy-hooker /usr/bin/
COPY config/template.tpl /etc/nginx/
COPY config/services.d /etc/services.d

VOLUME /etc/nginx

EXPOSE 80
ENTRYPOINT ["/bin/s6-svscan", "/etc/services.d"]
CMD []
