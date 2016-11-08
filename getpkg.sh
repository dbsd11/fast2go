#!/bin/sh
  oldgo=$GOPATH
  GOPATH=$(cd `dirname $0`; pwd)

  if [ ! -e bin/gvt ];then
      go get 'github.com/FiloSottile/gvt'
      mkdir src/vendor
  fi

  if [ -e src/vendor/manifest ];then
      cd src
      ../bin/gvt restore -connections 3
      cd ..
  fi

#  go get ./..

  GOPATH=$oldgo

