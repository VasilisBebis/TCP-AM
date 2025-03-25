# Introduction

A basic protocol, running over TCP, that can perform calculations over a set of given numbers.
Created as a part of a Lab exercise for an "Internet Protocols" course.

# Implementation
The program is implemented in [Go](https://go.dev/) using the [net](https://pkg.go.dev/net@go1.24.1) package from the [standard library](https://pkg.go.dev/std)

# Use
## How to run the server
```console
$ go run src/server/server-main.go
```
## How to run the client
```console
$ go run src/client/client-main.go
```
## How to generate documentation
After installing [godoc](https://pkg.go.dev/golang.org/x/tools/cmd/godoc) run:
```console
$ godoc -http=:8000
```
then navigate to [http://localhost:8000/pkg/github.com/VasilisBebis/TCP-AM/pkg/](http://localhost:8000/pkg/github.com/VasilisBebis/TCP-AM/pkg/) to see the documentation for the server and client package
