#!/bin/bash
set -e

cmd="$@"
date_to_wait_until=$(date -d "+1mins" +%s)
echo "Waiting Habitica API to be up @ $SERVER"
until curl -s $SERVER > /dev/null 2>&1; do
  >&2 printf "."
  sleep 1
  if [ $date_to_wait_until  -lt $(date +%s) ]; then
    echo
    echo "Habitica API took too long to start"
    exit 1
  fi
done

echo
>&2 echo "Habitica API is up at starting running '$cmd'"
exec $cmd