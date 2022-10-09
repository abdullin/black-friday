#!/bin/bash
set -eu -o pipefail -o xtrace
TARGET=$1
GIT=$2
BRANCH=$3


  #git --work-tree=$TARGET --git-dir=$(pwd) checkout -f $BRANCH && cd $TARGET && lib/qa.sh >/dev/null 2>&1 &

git --work-tree=$TARGET --git-dir=$GIT checkout -f $BRANCH
cd $TARGET && /usr/local/go/bin/go run *.go && lib/send.py