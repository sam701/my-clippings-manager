#!/bin/bash

$GOPATH/bin/go-bindata $1 -prefix web web
go install .

