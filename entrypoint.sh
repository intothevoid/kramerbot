#!/bin/sh
set -e

if [ -z "$JWT_SECRET" ]; then
  JWT_SECRET=$(openssl rand -hex 32)
  echo "[kramerbot] JWT_SECRET not set — generated a new one for this session"
  export JWT_SECRET
fi

exec ./kramerbot "$@"
