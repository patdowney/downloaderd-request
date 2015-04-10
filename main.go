package main

import (
	"flag"
	"io"
	"log"
	"net/url"
	"os"

	"github.com/patdowney/downloaderd-common/http"
	"github.com/patdowney/downloaderd-request/api"
	"github.com/patdowney/downloaderd-request/download"
	dh "github.com/patdowney/downloaderd-request/http"
	"github.com/patdowney/downloaderd-request/local"
	//	"github.com/patdowney/downloaderd-request/rethinkdb"
)

// Config ...
type Config struct {
	ListenAddress   string
	RequestDataFile string

	DownloadServiceURL string
	AccessLogWriter    io.Writer
	ErrorLogWriter     io.Writer

	RethinkDBAddress string
}

// ConfigureLogging ...
func ConfigureLogging(config *Config) {
	log.SetOutput(config.ErrorLogWriter)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
}

// ParseArgs ...
func ParseArgs() *Config {
	c := &Config{}
	flag.StringVar(&c.ListenAddress, "http", "localhost:8090", "address to listen on")
	flag.StringVar(&c.RethinkDBAddress, "rethinkdb", "localhost:28015", "address to connect to")
	flag.StringVar(&c.DownloadServiceURL, "downloadurl", "http://localhost:8080/download/", "download agent service")
	flag.StringVar(&c.RequestDataFile, "requestdata", "requests.json", "request database file")
	flag.Parse()

	c.AccessLogWriter = os.Stdout
	c.ErrorLogWriter = os.Stderr

	return c
}

// CreateServer ...
func CreateServer(config *Config) {
	s := http.NewServer(&http.Config{ListenAddress: config.ListenAddress}, os.Stdout)

	requestStore, err := local.NewRequestStore(config.RequestDataFile)
	/*
		c := rethinkdb.Config{Address: config.RethinkDBAddress,
			MaxIdle:  10,
			MaxOpen:  20,
			Database: "Downloaderd"}

		requestStore, err := rethinkdb.NewRequestStore(c)
	*/
	if err != nil {
		log.Printf("init-request-store-error: %v", err)
	}

	linkResolver := api.NewLinkResolver(s.Router)
	linkResolver.DefaultScheme = "http"
	linkResolver.DefaultHost = config.ListenAddress

	downloadURL, _ := url.Parse(config.DownloadServiceURL)

	downloadClient, _ := download.NewHTTPClient(downloadURL)

	requestService := download.NewRequestService(requestStore, downloadClient)

	requestResource := dh.NewRequestResource(requestService, linkResolver)
	s.AddResource("/request", requestResource)

	err = s.ListenAndServe()
	log.Printf("init-listen-error: %v", err)
}

func main() {
	config := ParseArgs()

	ConfigureLogging(config)

	CreateServer(config)
}
