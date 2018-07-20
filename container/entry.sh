#!/bin/ash

while [ true ];do
  ping ${REMOTE_HOST}
  sleep 1
done
