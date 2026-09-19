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
	username          string
	password          string
	listenAddress     string
	ingestion         string
	pollInterval      time.Duration
	reconcileInterval time.Duration
	showVersion       bool
)

// Credentials must not be flag defaults: flag prints defaults in its usage
// output, which leaks the password on -help or any flag error. They are read
// from the environment after parsing instead, so an explicit flag still wins.
func parseFlags() {
	flag.StringVar(&username, "username", "", "Easee username (defaults to $EASEE_USERNAME)")
	flag.StringVar(&password, "password", "", "Easee password (defaults to $EASEE_PASSWORD)")
	flag.StringVar(&listenAddress, "listen-address", ":8080", "Address to listen on for HTTP requests (defaults to :8080)")
	flag.StringVar(&ingestion, "ingestion", "", "Where charger state comes from: stream or poll (defaults to $EASEE_INGESTION, else stream)")
	flag.DurationVar(&pollInterval, "poll-interval", time.Minute, "How often to poll charger state with -ingestion=poll")
	flag.DurationVar(&reconcileInterval, "reconcile-interval", 15*time.Minute, "How often to reconcile streamed state with a poll; 0 disables")
	flag.BoolVar(&showVersion, "version", false, "Print version information and exit")
	flag.Parse()
	if username == "" {
		username = os.Getenv("EASEE_USERNAME")
	}
	if password == "" {
		password = os.Getenv("EASEE_PASSWORD")
	}
	if ingestion == "" {
		ingestion = os.Getenv("EASEE_INGESTION")
	}
	if ingestion == "" {
		ingestion = "stream"
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

// pollAll refreshes every charger over REST: one request each.
func pollAll(client *easee.Client, chargers []easee.Charger, collectors map[string]*metrics.ChargerStateCollector) {
	for _, c := range chargers {
		updateChargerState(client, c, collectors[c.Id])
	}
}

func pollLoop(ctx context.Context, client *easee.Client, chargers []easee.Charger, collectors map[string]*metrics.ChargerStateCollector, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			pollAll(client, chargers, collectors)
		}
	}
}

func main() {
	parseFlags()
	if showVersion {
		fmt.Printf("%s\n", version.Print("easee-exporter"))
		os.Exit(0)
	}
	if ingestion != "stream" && ingestion != "poll" {
		exit("Unknown -ingestion %q, want stream or poll\n", ingestion)
	}
	log.Printf("Starting easee-exporter %s\n", version.Version)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := easee.NewClient(ctx, username, password)
	chargers, err := client.Chargers()
	if err != nil {
		exit("Chargers request failed: %v\n", err)
	}
	log.Printf("Chargers: %v\n", chargers)

	collectors := make(map[string]*metrics.ChargerStateCollector, len(chargers))
	ids := make([]string, 0, len(chargers))
	for _, c := range chargers {
		log.Printf("Starting monitoring charger %v\n", c.Id)
		collector := metrics.NewChargerStateCollector(c.Id)
		prometheus.MustRegister(collector)
		collectors[c.Id] = collector
		ids = append(ids, c.Id)
	}

	// Seed over REST either way, so the first scrape is complete even if the
	// stream cannot be established at all.
	pollAll(client, chargers, collectors)

	if ingestion == "poll" {
		log.Printf("Polling charger state every %v\n", pollInterval)
		go pollLoop(ctx, client, chargers, collectors, pollInterval)
	} else {
		stream := easee.NewStream(client, ids, func(charger string, o *easee.Observation) {
			if collector, ok := collectors[charger]; ok {
				collector.ApplyObservation(o)
			}
		})
		prometheus.MustRegister(metrics.NewStreamCollector(stream))
		go func() {
			if err := stream.Run(ctx); err != nil && ctx.Err() == nil {
				log.Printf("Stream stopped: %v\n", err)
			}
		}()
		// A stream that is connected but silent leaves every metric frozen
		// at its last value with nothing to show for it, so a slow poll runs
		// alongside to correct anything the stream misses.
		if reconcileInterval > 0 {
			log.Printf("Streaming charger state, reconciling every %v\n", reconcileInterval)
			go pollLoop(ctx, client, chargers, collectors, reconcileInterval)
		} else {
			log.Printf("Streaming charger state, reconciliation disabled\n")
		}
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
