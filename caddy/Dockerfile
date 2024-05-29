FROM caddy:2.8-alpine

RUN apk add --no-cache gettext

COPY Caddyfile /etc/caddy/Caddyfile.template

COPY docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh

RUN chmod +x /usr/local/bin/docker-entrypoint.sh

EXPOSE 443

ENTRYPOINT ["docker-entrypoint.sh"]
