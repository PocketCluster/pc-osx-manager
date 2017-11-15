#!/usr/bin/env bash

GOOS=linux go build -ldflags="-s -w -v=2" ./...
upx --brute ./pocketd