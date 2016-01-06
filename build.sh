#!/usr/bin/env bash
set -x

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o docker2graphite.linux docker2graphite.go
