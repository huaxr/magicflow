#!/usr/bin/env bash


docker run --name myZookeeper -e JVMFLAGS="-Xmx1024m" -p 2181:2181 zookeeper
