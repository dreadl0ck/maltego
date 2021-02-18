# Maltego

[![Go Report Card](https://goreportcard.com/badge/github.com/dreadl0ck/maltego)](https://goreportcard.com/report/github.com/dreadl0ck/maltego)
[![License](https://img.shields.io/badge/license-MIT-green)](https://raw.githubusercontent.com/dreadl0ck/maltego/master/LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/dreadl0ck/maltego.svg)](https://pkg.go.dev/github.com/dreadl0ck/maltego)

This is a Go package that provides datastructures for interacting with the [Maltego](https://www.maltego.com) graphical link analysis tool.

## Key Features

- ***type safe datastructures for code components***: messages and configuration entities are both fully modeled.
- ***utility functions***: Escaping XML, Calculating line thickness etc.
- ***automatically escapes input***: user supplied values are guaranteed to be properly escaped so they don't break the XML.
- ***go modules support***: reproducible builds, semantic versioning
- ***unit tests***: functional correctness 
- ***usage examples***: well documented and containerized
- ***MIT licensed***: can be incorporated into proprietary products

## Installation

Install the library for use in your Go application:

    go get github.com/dreadl0ck/maltego

## Usage Examples

Check the examples folder and unit tests!

## Unit Tests

Run the unit tests

    go test ./...

## Code Stats

```shell
$ cloc *.go
      14 text files.
      14 unique files.                              
       0 files ignored.

github.com/AlDanial/cloc v 1.84  T=0.02 s (872.4 files/s, 156290.9 lines/s)
-------------------------------------------------------------------------------
Language                     files          blank        comment           code
-------------------------------------------------------------------------------
Go                              14            366            355           1787
-------------------------------------------------------------------------------
SUM:                            14            366            355           1787
-------------------------------------------------------------------------------
```

## License

MIT
