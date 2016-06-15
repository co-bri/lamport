# Contributing to Lamport

We welcome, and encourage, developers of all skill-levels to contribute to Lamport. Please do not hesitate to contribute because of fear of public shame, embarassment, or making mistakes. We're interested in providing a safe and constructive project for everyone interested in learning about distributed systems.

## Developmer Setup

Lamport is currently written in [Go](https://golang.org/), and you'll want to take a look at [How to Write Go Code](https://golang.org/doc/code.html) for local development instructions. If you are looking to get started with Go development, [A Tour of Go](https://tour.golang.org/welcome/1) is a great place to start

## Guidelines

The goal of Lamport is to encourage people to explore a working distributed system build from the ground up. That being said, we do have some basic guidelines for development:

* Make sure to format your code accroding to `go fmt`
* Run your changes through `go lint`, and `go vet`
* Make sure to provide tests that cover the feature you are adding

The last item is particularly important because verifying a distributed system operates correctly is crucial to the correctness of a distributed system.

## Issues

Future enhancements are added to the issue section of the project as discussed on the Lamport Developers Google Group. These issues are labeled as "enhancement", and "help wanted". This is the best place to start if you'd like to make a contribution.
