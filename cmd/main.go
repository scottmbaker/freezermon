package main

// freezermon
// Scott Baker
//
// Prometheus collector for making sure freezer doesn't go bad.
// Requires a DS18B20 temperature sensor (usually attached to GPIO4) and the 1w
// driver loaded in /boot/config.txt.
//
// To test:
// curl -H localhost:8080/metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/scottmbaker/freezermon/pkg/ds18b20"
	"github.com/sevlyar/go-daemon"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	DEVICE_DIR = "/sys/bus/w1/devices"
	PIDFILE    = "/tmp/freezermon.pid"
	LOGFILE    = "/tmp/freezermon.log"
	PROGNAME   = "freezermon"
)

var (
	verbose     bool
	runAsDaemon bool
	rootCmd     = &cobra.Command{
		Use:   PROGNAME,
		Short: "Freezer Monitor",
	}
)

func run() {
	// Create a DS18B20 object to collect the temperatures
	ds, err := ds18b20.NewDS18B20(verbose)
	if err != nil {
		log.Fatalf("Error initializing DS18B20: %v", err)
	}

	// Perform an initial reading -- if we don't succeed, fail fast.
	tempC, err := ds.MeasureFirstDevice()
	if err != nil {
		log.Fatalf("Error measuring initial temperature: %v", err)
	} else {
		log.Printf("Initial temperature: %.2f°C", tempC)
	}

	// Create the prometheus gauge
	temperature := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "temperature_celsius",
		Help: "Current temperature in Celsius",
	})
	up := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "up",
		Help: "Whether the collector is healthy (1 = healthy, 0 = failed)",
	})
	registry := prometheus.NewRegistry()
	registry.MustRegister(temperature)
	registry.MustRegister(up)

	// Create a goroutine to periodically measure the temperature
	go func() {
		// if for any reason the goroutine exits, make sure we set the up gauge to 0
		defer func() {
			fmt.Fprintf(os.Stderr, "Collection goroutine as exited, setting Up to 0\n")
			up.Set(0)
		}()
		for {
			tempC, err := ds.MeasureFirstDevice()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error measuring temperature: %v\n", err)
				temperature.Set(0) // Set to 0 on error
				up.Set(0)
			} else {
				if verbose {
					log.Printf("Sampled Temperature: %.2f°C", tempC)
				}
				temperature.Set(tempC)
				up.Set(1)
			}
			time.Sleep(10 * time.Second)
		}
	}()

	// Add the metrics handler to the http server
	http.Handle(
		"/metrics", promhttp.HandlerFor(
			registry,
			promhttp.HandlerOpts{
				EnableOpenMetrics: true,
			}),
	)

	// Start the HTTP server
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error starting HTTP server: %v\n", err)
		os.Exit(-1)
	}
}

// much thanks to sevlyar for the go-daemon library!

func daemonize() {
	cntxt := &daemon.Context{
		PidFileName: PIDFILE,
		PidFilePerm: 0644,
		LogFileName: LOGFILE,
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
	}

	d, err := cntxt.Reborn()
	if err != nil {
		log.Fatal("Unable to run: ", err)
	}
	if d != nil {
		return
	}
	defer cntxt.Release() // nolint:errcheck

	log.Print("- - - - - - - - - - - - - - -")
	log.Print("daemon started")

	run()
}

func mainCommand(cmd *cobra.Command, args []string) {
	if runAsDaemon {
		daemonize()
	} else {
		run()
	}
}

func main() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose debug messages")
	rootCmd.PersistentFlags().BoolVarP(&runAsDaemon, "daemon", "D", false, "run as a daemon")
	rootCmd.Run = mainCommand

	err := rootCmd.Execute()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
