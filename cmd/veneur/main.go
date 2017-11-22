package main

import (
	"flag"
	"os"
	"time"

	"github.com/getsentry/raven-go"
	"github.com/sirupsen/logrus"
	"github.com/stripe/veneur"
	"github.com/stripe/veneur/trace"
)

var (
	configFile = flag.String("f", "", "The config file to read for settings.")
	logFormat = flag.String("logFormat", "text", "Set log format \"text\" or \"JSON\"")
)

func init() {
	trace.Service = "veneur"
}

func main() {
	flag.Parse()

	if logFormat != nil {
		switch *logFormat {
		case "json":
			logrus.SetFormatter(&logrus.JSONFormatter{})
		case "text":
			logrus.SetFormatter(&logrus.TextFormatter{})
		default:
			logrus.Fatalf("Unknown log format %s", *logFormat)
		}
	}

	if configFile == nil || *configFile == "" {
		logrus.Fatal("You must specify a config file")
	}

	conf, err := veneur.ReadConfig(*configFile)
	if err != nil {
		logrus.WithError(err).Fatal("Error reading config file")
	}

	server, err := veneur.NewFromConfig(conf)
	if err != nil {
		e := err

		logrus.WithError(e).Error("Error initializing server")
		var sentry *raven.Client
		if conf.SentryDsn != "" {
			sentry, err = raven.New(conf.SentryDsn)
			if err != nil {
				logrus.WithError(err).Error("Error initializing Sentry client")
			}
		}

		hostname, _ := os.Hostname()

		p := raven.NewPacket(e.Error())
		if hostname != "" {
			p.ServerName = hostname
		}

		_, ch := sentry.Capture(p, nil)
		select {
		case <-ch:
		case <-time.After(10 * time.Second):
		}

		logrus.WithError(e).Fatal("Could not initialize server")
	}
	defer func() {
		veneur.ConsumePanic(server.Sentry, server.Statsd, server.Hostname, recover())
	}()

	if server.TraceClient != nil {
		if trace.DefaultClient != nil {
			trace.DefaultClient.Close()
		}
		trace.DefaultClient = server.TraceClient
	}
	server.Start()

	if conf.HTTPAddress != "" {
		server.HTTPServe()
	} else {
		select {}
	}
}
