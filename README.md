# Lamport

An academic exercise in building a distributed system, named in honor of Turing Award Winner [Leslie Lamport](http://www.lamport.org/). The purpose of this project is to explore the complexities involved with building a distributed computing system. The origins of this project are a [2016 presentation](http://www.meetup.com/Distributed-Computing-Denver/events/230054258/) for the [Distributed Computing Denver](http://www.meetup.com/Distributed-Computing-Denver/) Meetup group.

## Developer Setup

To get your local development environment setup, follow these steps:

1. Download [Git](https://git-scm.com/downloads) and follow the [first time Git setup](https://git-scm.com/book/en/v2/Getting-Started-First-Time-Git-Setup)
2. Take a moment to walkthrough the [Getting Started](https://golang.org/doc/install) guide to install Go to your local development environment
3. Read through the ["How to Write Go Code"](https://golang.org/doc/code.html) to setup your GOPATH and workspace
4. Clone this repository: `git clone https://github.com/Distributed-Computing-Denver/lamport.git`

## Running

If you have the [Make](https://www.gnu.org/software/make/) utility installed, you can build a Lamport artifact by running:

`make all`

If not, run through the following steps.

Once you've cloned the repository as outlined above, and made sure to [setup your workspace](https://golang.org/doc/code.html), navigate to the root folder of the project and run:

`go build`

Assuming you've setup your Go workspace correctly, an exectuable file named "lamport" is created. You can now start lamport by running:

`./lamport run`

You can get information about Lamport's command line flags by running:

`./lamport -h`

Coming soon will be build scripts and infrastructure as code that will allow you to get Lamport and it's depedencies up and running (pull requests are very much welcome). The idea is that operations folks can run Lamport to get a feel for running a distributed system. Lamport intends to be fully operationalized as each feature is added. Stay tuned for more details.

## Configuration 

Lamport is configured via a [TOML](https://github.com/toml-lang/toml) configuration file. The application defaults to lamport.toml, though you can specify a different configuration file via a command line arg. The following values are defined in the file:

- **Host** - The hostname that Lamport runs on and uses to advertise connections. The default is 127.0.0.1
- **Port** - The port name that the Lamport server uses to connect. The default is 5936.

We're also working on a Wiki that will walk users through the implementation details for specific features. Stay tuned for more details.

## Testing

To run tests:

`go test ./...`

## License

MIT
