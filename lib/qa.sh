#!/bin/bash
set -eu -o pipefail
#-o xtrace
make specs
#go run *.go
#lib/send.py