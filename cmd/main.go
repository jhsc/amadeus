package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	logger       = log.New(os.Stdout, "", log.LstdFlags|log.LUTC)
	useEnvConfig = flag.Bool("e", false, "use environment variables as config")
)

func main() {
	flag.Usage = help
	flag.Parse()

	commands := map[string]func(){
		"start":   startServer,
		"init":    initConfig,
		"gen-key": genKey,
		"help":    help,
	}

	if cmdFunc, ok := commands[flag.Arg(0)]; ok {
		cmdFunc()
	} else {
		help()
		//  Misuse of shell builtins
		os.Exit(1)
	}
}

func startServer() {
	logger.Println("start server srv")
}

func initConfig() {
	logger.Println("initiliaze config srv")
}

func genKey() {
	logger.Println("gen keys srv")
}

func help() {
	fmt.Fprintln(os.Stderr, `Usage:
	amadeus start                      - start the server
	amadeus init                       - create an initial configuration file
	amadeus gen-key                    - generate a random 32-byte hex-encoded key
	amadeus help                       - show this message
Use -e flag to read configuration from environment variables instead of a file. E.g.:
	amadeus -e start`)
}
