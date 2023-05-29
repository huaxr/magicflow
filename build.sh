#!/usr/bin/env bash
export GOPROXY=goproxy.io,goproxy.xxx.com

go env -w GOPRIVATE="github.com"
cat  ~/.gitconfig
make api
make rpc
