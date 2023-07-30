#! /usr/bin/env bash

#Avoid (cached) results
go test -v . -test.count=1 \
  -test.parallel 5 -timeout 30s

go test -bench=..
# go test -cover -coverprofile=c.out && go tool cover -html=c.out -o coverage.html