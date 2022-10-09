#!/bin/bash
# git --work-tree=$TARGET --git-dir=$DIR checkout -f $BRANCH >> /dev/null
/usr/local/go/bin/go run *.go && lib/send.py