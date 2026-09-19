package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	promver "github.com/prometheus/client_golang/prometheus/collectors/version"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	"github.com/terjesannum/easee-exporter/internal/easee"
	"github.com/terjesannum/easee-exporter/internal/metrics"
)

var (
	username      string
	password      string
	listenAddress string
	showVersion   bool
)

// Credentials must not be flag defaults: flag prints defaults in its usage
// output, which leaks the password on -help or any flag error. They are read
// from the environment after parsing instead, so an explicit flag still wins.
func parseFlags() {
	flag.StringVar(&username, "username", "", "Easee username (defaults to $EASEE_USERNAME)")
	flag.StringVar(&password, "password", "", "Easee password (defaults to $EASEE_PASSWORD)")
	flag.StringVar(&listenAddress, "listen-address", ":8080", "Address to listen on for HTTP requests (defaults to :8080)")
	flag.BoolVar(&showVersion, "version", false, "Print version information and exit")
	flag.Parse()
	if username == "" {
		username = os.Getenv("EASEE_USERNAME")
	}
	if password == "" {
		password = os.Getenv("EASEE_PASSWORD")
	}
}

func exit(format string, v ...any) {
	log.Printf(format, v...)
	time.Sleep(time.Second * 10)
	os.Exit(1)
}

func updateChargerState(client *easee.Client, charger easee.Charger, collector *metrics.ChargerStateCollector) {
	state, err := client.ChargerState(charger.Id)
	if err != nil {
		// Keep serving the last known state: exiting here restarts the
		// process, and the repeated logins burn through the API rate limit.
		log.Printf("Charger state request failed for %s: %v\n", charger.Id, err)
		return
	}
	collector.UpdateState(&state)
	log.Printf("Updated charger state for %s\n", charger.Id)
}

func main() {
	parseFlags()
	if showVersion {
		fmt.Printf("%s\n", version.Print("easee-exporter"))
		os.Exit(0)
	}
	log.Printf("Starting easee-exporter %s\n", version.Version)
	ctx := context.Background()
	client := easee.NewClient(ctx, username, password)
	chargers, err := client.Chargers()
	if err != nil {
		exit("Chargers request failed: %v\n", err)
	}
	log.Printf("Chargers: %v\n", chargers)
	for _, c := range chargers {
		log.Printf("Starting monitoring charger %v\n", c.Id)
		collector := metrics.NewChargerStateCollector(c.Id)
		prometheus.MustRegister(collector)
		updateChargerState(client, c, collector)
		ticker := time.NewTicker(time.Minute)
		quit := make(chan struct{})
		go func(charger easee.Charger) {
			for {
				select {
				case <-ticker.C:
					updateChargerState(client, charger, collector)
				case <-quit:
					ticker.Stop()
					return
				}
			}
		}(c)
	}
	prometheus.MustRegister(promver.NewCollector("easee_exporter"))
	log.Println("Starting http listener")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Easee prometheus exporter")
	})
	http.Handle("/metrics", promhttp.Handler())
	err = http.ListenAndServe(listenAddress, nil)
	exit("Error: %v", err)
}
