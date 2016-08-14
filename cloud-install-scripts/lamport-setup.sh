#!/bin/bash
#
# Copyright (c) 2016 Brian McKean
#
# Assume new intall on ubuntu
# Do all in shell script now -- later set up for ansible

set -x
# Get any updates
sudo apt-get -y update

# Setup for go
# Note: apt-get for go is too old version  --- must do manual isntall of go
cd 
sudo apt-get -y install curl
curl -O https://storage.googleapis.com/golang/go1.6.linux-amd64.tar.gz
tar -xvf go1.6.linux-amd64.tar.gz
sudo cp -r go /usr/local
sudo ln -s /usr/local/go/bin/go /usr/bin/go

# Make directories for go
mkdir work
cd work
mkdir src
mkdir pkg
mkdir bin

#Setup to use something built in go
export GOPATH=$HOME/work
export PATH=$PATH:$GOPATH/bin

#install git
sudo apt-get -y install git

#install golint
go get -u github.com/golang/lint/golint

#Install make
sudo apt-get -y install build-essential

# Get lamport code
cd src
mkdir github.com
cd github.com
mkdir Distributed-Computing-Denver
cd Distributed-Computing-Denver

git init
git clone https://github.com/Distributed-Computing-Denver/lamport.git

cd lamport

# Build lamport
# Note: not sure why make needed --- but it downloads some 
# other code from github that make all does not download
make
make all

# fixme
# add port specific setup
#

