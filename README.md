# Amadeus

Deployer tool

## Installation

Amadeus requires Go 1.9 or later.

```sh
go get github.com/jhsc/amadeus
```

## Usage

TODO

## Development

### Requirements

- Install [Go](https://golang.org)
- Install [Go Modules](https://blog.golang.org/using-go-modules)

### Makefile

```sh
// Clean up
$ make clean

// Creates folders and download dependencies
$ make configure

//Run tests and generates html coverage file
make cover

// Download project dependencies
make depend

// Format all go files
make fmt

//Run linters
make lint

// Run tests
make test
```

## License

This project is released under the MIT licence. See [LICENSE](https://github.com/jhsc/amadeus/blob/master/LICENSE) for more details.
