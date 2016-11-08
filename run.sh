#!/bin/sh
  GOPATH=$GOPATH:$(cd `dirname $0`; pwd)
#  cd src

  go run $1

 # cd ..