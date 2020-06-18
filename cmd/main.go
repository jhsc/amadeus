package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"gitlab.com/jhsc/amadeus/api"
	"gitlab.com/jhsc/amadeus/config"
	"gitlab.com/jhsc/amadeus/docker"
	"gitlab.com/jhsc/amadeus/store/sqlite"
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
	cfg, err := config.GetConfig(*useEnvConfig)
	if err != nil {
		logger.Fatalf("failed to load configuration: %s", err)
	}
	baseURL, err := url.Parse(cfg.BaseURL)
	if err != nil {
		logger.Fatalf("failed to parse base url: %s", err)
	}

	store, err := sqlite.Connect()
	if err != nil {
		logger.Fatalf("failed to create new docker service: %s", err)
	}

	endpoint := "unix:///var/run/docker.sock"
	ds, err := docker.NewService(endpoint, store)
	if err != nil {
		logger.Fatalf("failed to create new docker service: %s", err)
	}

	apiHandler := api.New(&api.Config{
		Logger:        logger,
		DockerService: ds,
		Store:         store,
	},
		cfg.Token,
	)

	router := chi.NewRouter()
	router.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: logger}))
	router.Use(middleware.Recoverer)

	router.Get("/healthz", apiHandler.HandleHealthz)
	router.Mount("/api/v1", apiHandler)

	if err := http.ListenAndServe(cfg.Address, http.StripPrefix(baseURL.Path, router)); err != nil {
		logger.Fatalf("listen and serve failed: %v", err)
	}
}

func initConfig() {
	if _, err := os.Stat(config.ConfigFile); !os.IsNotExist(err) {
		logger.Fatalf("configuration file already exists: %s", config.ConfigFile)
	}

	logger.Printf("creating initial configuration: %s", config.ConfigFile)

	cfg, err := config.Init()
	if err != nil {
		logger.Fatalf("failed to generate initial configuration: %s", err)
	}

	err = ioutil.WriteFile(config.ConfigFile, []byte(cfg), 0666)
	if err != nil {
		logger.Fatalf("failed to write configuration file: %s", err)
	}
}

func genKey() {
	logger.Printf("key: %s", config.GenKeyHex(32))
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
