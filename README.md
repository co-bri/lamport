# Lamport

An academic exercise in building a distributed system, named in honor of Turing Award Winner [Leslie Lamport](http://www.lamport.org/). The purpose of this project is to explore the complexities involved with building a distributed computing system. The origins of this project are a [2016 presentation](http://www.meetup.com/Distributed-Computing-Denver/events/230054258/) for the [Distributed Computing Denver](http://www.meetup.com/Distributed-Computing-Denver/) Meetup group.

## Developer Setup

To get your local development environment setup, follow these steps:

1. Download [Git](https://git-scm.com/downloads) and follow the [first time Git setup](https://git-scm.com/book/en/v2/Getting-Started-First-Time-Git-Setup)
2. Take a moment to walkthrough the [Getting Started](https://golang.org/doc/install) to install Go into your local development environment
3. Read through the ["How to Write Go Code"](https://golang.org/doc/code.html) to setup your GOPATH and workspace
4. Clone this repository: `git clone https://github.com/Distributed-Computing-Denver/lamport.git`

Lamport uses Git tags to provide developers with a way to see incremental features being added to the project. A new developer will want to fetch the lowest numbered tag/feature. The wiki will be updated to provide users context on the feature that was added, as well as some design decisions that imfluenced the choices.

## Running

Once you've cloned the repository as outlined above, and made sure to [setup your workspace](https://golang.org/doc/code.html), navigate to the root folder of the project and run:

`go build`

Assuming you've setup your Go workspace correctly, an exectuable file named "lamport" is created. You can now start lamport by running:

`./lamport`

Coming soon will be build scripts and infrastructure as code that will allow you to get Lamport and it's depedencies up and running (pull requests are very much welcome). The idea is that operations folks can run Lamport to get a feel for running a distributed system. Lamport intends to be fully operationalized as each feature is added. Stay tuned for more details.

## Testing

To run tests for lamport you'll need to install [Apache Zookeeper](https://zookeeper.apache.org/releases.html) locally, and add the 'bin' directory of the distribution to the $PATH variable of the user running the tests. One this is done, you can execute tests by running the following command:

`go test ./...`

The tests will spin up a local instance of Zookeeper, and clean up any changes made to Zookeeper within the test suite.

### Zookeeper Install: Linux

Users running Linux can install Zookeeper by running the following commands:

```
wget http://apache.claz.org/zookeeper/zookeeper-3.4.8/zookeeper-3.4.8.tar.gz && 
tar -xvf zookeeper-3.4.8.tar.gz -C /usr/local
```

This will install Zookeeper to `/usr/local/zookeeper-3.4.8.tar.gz`, and then you can add `/usr/local/zookeeper-3.4.8/bin` to your path:

```
export PATH="$PATH:/usr/local/zookeeper-2.4.8/bin"
```

### Zookeeper Install: OSX

Users running OSX can install and manage ZooKeeper using [Homebrew](http://brew.sh/) by running:

```
brew install zookeeper
```

This will install ZooKeeper to `/usr/local/Cellar/zookeeper/X.Y.Z/` (this will vary depending on the current version of ZooKeeper, so make sure to note the version numbers). Add this to your path:

```
export PATH="$PATH:/usr/local/Cellar/zookeeper/X.Y.Z/bin"
```

Finally, create a symbolic link:

```
ln -s /usr/local/Cellar/zookeeper/X.Y.Z/bin/zkServer /usr/local/Cellar/zookeeper/X.Y.Z/bin/zkServer.sh
```

You can now run the test suite with:

```
go test ./...
```

## License

MIT
