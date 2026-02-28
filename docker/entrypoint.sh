#!/bin/sh
set -e

# Start nginx in background (serves frontend + reverse proxy)
nginx

# Start Go backend in foreground
exec ./tinyrdm-server
