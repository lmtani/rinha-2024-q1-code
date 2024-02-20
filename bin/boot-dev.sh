#!/usr/bin/env sh

# autowatch and run tests
find . | entr -nr go test $(go list ./...) &

# autowatch and run server
find . | entr -nr go run cmd/server/*.go
