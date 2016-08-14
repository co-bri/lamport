#
# Copyright (c) 2016 Brian McKean
#
# using ssh on remote server ...

##Run basic update
sudo apt-get update

##Install Git
sudo apt-get -y install git


## Git Setup
git config --global user.name "Brian McKean"
git config --global user.email bdmckean@gmail.com


## Go Setup
sudo apt-get Ðy install gccgo-go

## install zookeeper
wget http://apache.claz.org/zookeeper/zookeeper-3.4.8/zookeeper-3.4.8.tar.gz
sudo tar -xvf zookeeper-3.4.8.tar.gz -C /usr/local
export PATH="$PATH:/usr/local/zookeeper-2.4.8/bin"

## Set up lamport & go build area
mkdir lamport_proj
cd lamport_proj

mkdir src
mkdir pkg
mkdir bin

cd src

export GOPATH=$HOME/lamport_proj

mkdir github.com
cd github.com
mkdir Distributed-Computing-Denver
cd Distributed-Computing-Denver

/* Why do we have all these directories? */

## get code
git init
git clone https://github.com/Distributed-Computing-Denver/lamport.git

cd lamport

## try to build lamport
go build

