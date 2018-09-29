package main

import (
	"io/ioutil"
	"os"
	"fmt"
	"net/http"
	"gopkg.in/yaml.v2"
	"github.com/jessevdk/go-flags"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"time"
)

const (
	Author  = "webdevops.io"
	Version = "0.1.0"
	AZURE_RESOURCEGROUP_TAG_PREFIX = "tag_"
)

var (
	argparser          *flags.Parser
	args               []string
	Logger             *DaemonLogger
	ErrorLogger        *DaemonLogger
)

var opts struct {
	// general settings
	Verbose     []bool `       long:"verbose" short:"v"  env:"VERBOSE"       description:"Verbose mode"`

	// server settings
	ServerBind  string `       long:"bind"               env:"SERVER_BIND"   description:"Server address"               default:":8080"`
	ScrapeTime  time.Duration `long:"scrape-time"        env:"SCRAPE_TIME"   description:"Scrape time (time.duration)"  default:"5s"`

	ConfigurationFile  string `long:"config-file"        env:"CONFIG"        description:"Configuration file"           required:"true"`
	configuration      Configuration
}

func main() {
	initArgparser()

	Logger = CreateDaemonLogger(0)
	ErrorLogger = CreateDaemonErrorLogger(0)

	// set verbosity
	Verbose = len(opts.Verbose) >= 1

	Logger.Messsage("Init Azure DevOps exporter v%s (written by %v)", Version, Author)

	Logger.Messsage("Init configuration")
	initConfiguration()

	Logger.Messsage("Starting metrics collection")
	setupMetricsCollection()
	startMetricsCollection()

	Logger.Messsage("Starting http server on %s", opts.ServerBind)
	startHttpServer()
}

// init argparser and parse/validate arguments
func initArgparser() {
	argparser = flags.NewParser(&opts, flags.Default)
	_, err := argparser.Parse()

	// check if there is an parse error
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			fmt.Println(err)
			fmt.Println()
			argparser.WriteHelp(os.Stdout)
			os.Exit(1)
		}
	}

	if _, err := os.Stat(opts.ConfigurationFile); os.IsNotExist(err) {
		fmt.Println(fmt.Sprintf("configuration file \"%v\" doesn't exists", opts.ConfigurationFile))
		fmt.Println()
		argparser.WriteHelp(os.Stdout)
		os.Exit(1)
	}
}

// Init and build Azure authorzier
func initConfiguration() {
	data, err := ioutil.ReadFile(opts.ConfigurationFile)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal([]byte(data), &opts.configuration)
	if err != nil {
		panic(err)
	}
}

// start and handle prometheus handler
func startHttpServer() {
	http.Handle("/metrics", promhttp.Handler())
	ErrorLogger.Fatal(http.ListenAndServe(opts.ServerBind, nil))
}
