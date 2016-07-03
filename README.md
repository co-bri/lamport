# Lamport

An academic exercise in building a distributed system, named in honor of Turing Award Winner [Leslie Lamport](http://www.lamport.org/). The purpose of this project is to explore the complexities involved with building a distributed computing system. The origins of this project are a [2016 presentation](http://www.meetup.com/Distributed-Computing-Denver/events/230054258/) for the [Distributed Computing Denver](http://www.meetup.com/Distributed-Computing-Denver/) Meetup group.

## Developer Setup

To get your local development environment setup, follow these steps:

1. Download [Git](https://git-scm.com/downloads) and follow the [first time Git setup](https://git-scm.com/book/en/v2/Getting-Started-First-Time-Git-Setup)
2. Take a moment to walkthrough the [Getting Started](https://golang.org/doc/install) guide to install Go to your local development environment
3. Read through the ["How to Write Go Code"](https://golang.org/doc/code.html) to setup your GOPATH and workspace
4. Clone this repository: `git clone https://github.com/Distributed-Computing-Denver/lamport.git`

## Running

Once you've cloned the repository as outlined above, and made sure to [setup your workspace](https://golang.org/doc/code.html), navigate to the root folder of the project and run:

`go build`

Assuming you've setup your Go workspace correctly, an exectuable file named "lamport" is created. You can now start lamport by running:

`./lamport`

Please note that you need an instance of a Zookeeper server running if Lamport is using the Zookeeper library for leader election in order for Lamport to start successfully.

You can get information about Lamport's command line flags by running:

`./lamport -h`

Coming soon will be build scripts and infrastructure as code that will allow you to get Lamport and it's depedencies up and running (pull requests are very much welcome). The idea is that operations folks can run Lamport to get a feel for running a distributed system. Lamport intends to be fully operationalized as each feature is added. Stay tuned for more details.

## Configuration 

Lamport is configured via a [TOML](https://github.com/toml-lang/toml) configuration file. The application defaults to lamport.toml, though you can specify a different configuration file via a command line arg. The following values are defined in the file:

- **ElectionLibrary** - The library to use for leader elections. This value must be either "Raft" or "Zookeeper". The default is Zookeeper.
- **Host** - The hostname that Lamport runs on and uses to advertise connections. The default is 127.0.0.1
- **LamportPort** - The port name that the Lamport server uses to connect. The default is 5936.
- **RaftDir** - The directory the Raft library uses to store data about the state of the Raft cluster. The default is .raftDir.
- **RaftPort** - The port the Raft library uses to communicate on. The default is 8500.

We're also working on a Wiki that will walk users through the implementation details for specific features (e.g. leader election)

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
