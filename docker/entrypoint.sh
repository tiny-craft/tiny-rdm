#!/bin/sh
set -e

# Pass through environment variables
# PORT - HTTP listen port (default: 8088)
# REDIS_HOST - default Redis host for auto-connect (optional)
# REDIS_PORT - default Redis port (optional)
# REDIS_PASSWORD - default Redis password (optional)

exec "$@"
