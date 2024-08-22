#!/bin/bash

HOMEDIR=$(dirname $(dirname $(readlink -f $0)))
OUTDIR=${HOMEDIR}/output

cd ${HOMEDIR} && rm -rf ${OUTDIR}  2>/dev/null

# 编译mac下可以执行文件
#go build -ldflags "-s -w" -o commander-mac main.go

# 使用交叉编译 linux和windows版本可以执行的文件
CGO_ENABLED=0 GOOS=linux   GOARCH=amd64   go build -ldflags "-s -w" -o ${OUTDIR}/commander-linux main.go
CGO_ENABLED=0 GOOS=linux   GOARCH=arm64   go build -ldflags "-s -w" -o ${OUTDIR}/commander-linux-arm64 main.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64   go build -ldflags "-s -w" -o ${OUTDIR}/commander-win.exe main.go
CGO_ENABLED=0 GOOS=linux   GOARCH=loong64 go build -ldflags "-s -w" -o ${OUTDIR}/commander-linux-loong64 main.go
CGO_ENABLED=0 GOOS=darwin  GOARCH=arm64   go build -ldflags "-s -w" -o ${OUTDIR}/commander-mac main.go             # m1 芯片
CGO_ENABLED=0 GOOS=darwin  GOARCH=amd64   go build -ldflags "-s -w" -o ${OUTDIR}/commander-mac-intel64 main.go     # intel64