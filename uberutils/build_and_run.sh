#! /bin/bash
set -ex

make kube-apiserver
cp _output/bin/kube-apiserver "$GOPATH/bin/kube-apiserver"
bash uberutils/run-apiserver.sh 2>&1 | less
