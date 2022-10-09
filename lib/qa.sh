#!/bin/bash
set -eu -o pipefail
#-o xtrace
go run *.go
lib/send.py