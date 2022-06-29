#!/bin/bash
go build -tags release -ldflags "-w -s" -o dist/doss-dispatcher
build/upx -9 dist/doss-dispatcher
cp -r banner.txt config.yml db dist