#!/bin/sh

# Wait until etcd is ready.
while true; do
  if etcdctl endpoint health; then
    echo "etcd is ready."
    break
  fi
  echo "waiting for etcd to start..."
  sleep 1
done

# Publish the initial configuration.
etcdctl put /config/config.yaml -- < /config/config.yaml


# Watch for configuration changes.
previous_hash=$(sha256sum /config/config.yaml | awk '{print $1}')

while true; do
  current_hash=$(sha256sum /config/config.yaml | awk '{print $1}')

  if [ "$current_hash" != "$previous_hash" ]; then
    etcdctl put /config/config.yaml -- < /config/config.yaml
    echo "spot update, config updated in etcd. $(date +'%Y-%m-%d %H:%M:%S')"
    previous_hash="$current_hash"
  fi

  sleep 60
done
