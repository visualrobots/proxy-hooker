FROM alpine:3.4
RUN apk update && apk add nginx bash s6 && rm -rf /var/cache/apk/* && touch /etc/nginx/conf.d/vhosts.conf

COPY proxy-hooker /usr/bin/
COPY config/nginx.conf /etc/nginx/nginx.conf
COPY config/template.tpl /etc/nginx/
COPY config/services.d /etc/services.d

EXPOSE 80

ENTRYPOINT ["/bin/s6-svscan", "/etc/services.d"]
CMD []
