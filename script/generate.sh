#!/usr/bin/env bash

readonly exec_pos="$(dirname "$0")/../"

cd "$exec_pos" &&
  go get -u github.com/dmarkham/enumer &&
  go generate ./... &&
  go mod tidy
