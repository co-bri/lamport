#!/bin/bash
#
# Copyright (c) 2016 Brian McKean
#
export GOPATH=$HOME/work
export PATH=$PATH:$GOPATH

cd $HOME/work/src/github.com/Distributed-Computing-Denver/lamport

make all

wait
