#!/usr/bin/env bash
export GOPROXY=goproxy.io,goproxy.xesv5.com

go env -w GOPRIVATE="github.com"
cat  ~/.gitconfig
make api
make rpc
