# Amadeus

Deployer tool

## Installation

Amadeus requires Go 1.9 or later.

```sh
go get github.com/jhsc/amadeus
```

## Usage

Display help menu

```sh
amadeus help
Usage:
	amadeus start                      - start the server
	amadeus init                       - create an initial configuration file
	amadeus gen-key                    - generate a random 32-byte hex-encoded key
	amadeus help                       - show this message
Use -e flag to read configuration from environment variables instead of a file. E.g.:
	amadeus -e start
```

Generate Token

```sh
amadeus gen-key
2020/06/17 20:02:45 key: a0e95ec1c52387aa35fa885c50345755accfeccfb8e37f3bd5994d7d4198fc1e
```

Start server

```sh
amadeus start
```

Deployment example

```sh
curl --request POST \
  --url http://0.0.0.0:8080/api/v1/deploy \
  --header 'content-type: application/json' \
  --header 'token: TOKEN' \
  --data '{
  "project": "PROJECT",
  "compose_file": "BASE64 COMPOSE",
  "registry": {
    "url": "URL",
    "login": "USER",
    "password": "PASSWORD"
  },
  "extra": {
    "TAG": "0.0.0"
  }
}'
```

## Development

### Requirements

- Install [Go](https://golang.org)
- Go Modules [Go Modules](https://blog.golang.org/using-go-modules)

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
