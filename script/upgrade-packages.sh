#!/usr/bin/env bash

readonly exec_pos="$(dirname "$0")/../"

cd "$exec_pos" &&
  go get -d -u ./... && go mod tidy
