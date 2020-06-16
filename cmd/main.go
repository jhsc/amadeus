package main

import (
	"log"
	"os"
)

var (
	logger = log.New(os.Stdout, "", log.LstdFlags|log.LUTC)
)

func main() {
	logger.Println("hello :)")
}
