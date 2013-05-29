package main

import (
	"flag"
	"io/ioutil"
	dlog "log"
	"log/syslog"
	"os"

	"github.com/kilburn/gorbs/config"
	"github.com/kilburn/gorbs/log"

	"launchpad.net/goyaml"

	"github.com/nightlyone/lockfile"
)

var configPath = flag.String("config", "/etc/gorbs.conf", "Specify alternate config file (-c /path/to/file)")
var flagVerbose = flag.Bool("verbose", false, "Show equivalent shell commands being executed.")
var flagTest = flag.Bool("test", false, "Show verbose output, but don't touch anything. This will be similar, but not always exactly the same as the real output from a live run.")
var flagQuiet = flag.Bool("quiet", false, "Suppress non-fatal warnings.")
var flagVVerbose = flag.Bool("extra-verbose", false, "The same as -v, but with more detail.")
var flagDebug = flag.Bool("debug", false, "A firehose of diagnostic information.")

var conf *config.Configuration

func loadConfiguration() {
	reader, err := os.Open(*configPath)
	if err != nil {
		log.Alertf("Unable to read the configuration file (%v)", err.Error())
	}
	data, err := ioutil.ReadAll(reader)

	goyaml.Unmarshal(data, conf)
	if err != nil {
		log.Alert(err.Error())
	}
}

func setup() {
	flag.Parse()

	var verbose log.Level
	switch {
	case *flagDebug:
		verbose = log.DEBUG
	case *flagVVerbose:
		verbose = log.INFO
	case *flagVerbose:
		verbose = log.WARN
	case *flagQuiet:
		verbose = log.ALERT
	default:
		verbose = log.ERROR
	}

	// Apply command-line verbosity first
	logger := dlog.New(os.Stdout, "", 0)
	backend := log.NewStdBackend(verbose, logger)
	log.AddBackend("std", backend)

	// Load the configuration file
	conf = config.New(verbose)
	loadConfiguration()

	// Enforce command-line verbosity over configuration
	if verbose != log.ERROR {
		conf.Verbose = verbose
	}

	setupLogging()

	// Enforce command-line testing over configuration
	conf.Test = *flagTest
}

func setupLogging() {
	// Setup the syslog backend
	writer, err := syslog.New(syslog.LOG_USER, "gobs")
	if err != nil {
		log.Alertf("Unable to initialize syslogging (%v)", err.Error())
	}
	syslogger := log.NewSyslogBackend(conf.Verbose, writer)
	log.AddBackend("syslog", syslogger)

	// Setup the std backend
	log.GetBackend("std").SetLevel(conf.Verbose)

	// Setup the logfile backend
	// @Todo
}

func foo() {
	log.Debug("This is debug!")
	log.Info("This is an info")
	log.Warn("This is a warning!")
	log.Error("This is an error!")
	log.Panic("This is a panic!")
}

func acquireLock() lockfile.Lockfile {
	lock, err := lockfile.New(conf.Lockfile)
	if err == lockfile.ErrNeedAbsPath {
		log.Panicf(
			`Error: you must specify an absolute path for the lockfile.
	("%s" is not an absolute path)`, conf.Lockfile)
	}
	if err != nil {
		log.Panicf("Something wrong happened (%v)", err)
	}

	err = lock.TryLock()
	if err == lockfile.ErrBusy {
		log.Panicf("Another instance of gorbs is already running. Aborting.")
	}
	if err != nil {
		log.Panicf("Unable to acquire the lockfile (%v). Aborting.", err)
	}

	return lock
}

func main() {
	// Catch-all to avoid dumping the stack on panic
	defer func() {
		if conf.Verbose < log.DEBUG {
			err := recover()
			if err != nil {
				os.Exit(1)
			}
		}
	}()

	setup()

	lock := acquireLock()
	defer lock.Unlock()

}
