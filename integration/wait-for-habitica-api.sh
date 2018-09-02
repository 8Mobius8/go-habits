#!/bin/bash

set -e

cmd="$@"
date_to_wait_until=$(date -d "+1mins" +%s)
until go-habits isup; do
  >&2 echo "Habitica API is down at $SERVER"
  sleep 1
  if [ $date_to_wait_until  -lt $(date +%s) ]; then
    echo "Habitica API took too long to start"
    exit 1
  fi
done

>&2 echo "Habitica API is up at starting running '$cmd'"
exec $cmd