#!/usr/bin/env sh

set -eu

envsubst '$APP_HOST' < /etc/caddy/Caddyfile.template > /etc/caddy/Caddyfile

caddy fmt --overwrite /etc/caddy/Caddyfile
caddy validate --config /etc/caddy/Caddyfile

caddy run --config /etc/caddy/Caddyfile --adapter caddyfile