package easee

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/philippseith/signalr"
)

// streamUrl is Easee's observation stream. Note the host: the hubs moved off
// api.easee.cloud to streams.easee.com.
const streamUrl = "https://streams.easee.com/hubs/chargers"

// Stream keeps charger state fresh from Easee's SignalR stream instead of
// polling. Every charger shares one websocket, so an account of any size
// costs a single connection rather than one request per charger per tick.
//
// Subscribing asks the server to replay the charger's current state before
// switching to push on change, so a fresh connection yields a full picture
// rather than only what happens to change next.
type Stream struct {
	client   *Client
	chargers []string
	apply    func(charger string, o *Observation)
	url      string

	connected  atomic.Bool
	reconnects atomic.Int64
}

func NewStream(client *Client, chargers []string, apply func(charger string, o *Observation)) *Stream {
	return &Stream{client: client, chargers: chargers, apply: apply, url: streamUrl}
}

func (s *Stream) Connected() bool   { return s.connected.Load() }
func (s *Stream) Reconnects() int64 { return s.reconnects.Load() }

// Run maintains the stream until the context is cancelled.
func (s *Stream) Run(ctx context.Context) error {
	client, err := signalr.NewClient(ctx,
		signalr.WithConnector(func() (signalr.Connection, error) { return s.connect(ctx) }),
		signalr.WithReceiver(&streamReceiver{stream: s}),
		signalr.Logger(signalrLogger{}, false),
		signalr.WithBackoff(func() backoff.BackOff {
			b := backoff.NewExponentialBackOff()
			// Never give up. The library default stops after fifteen
			// minutes, which would leave the exporter running with a dead
			// stream and silently stale metrics.
			b.MaxElapsedTime = 0
			b.MaxInterval = 5 * time.Minute
			return b
		}),
	)
	if err != nil {
		return fmt.Errorf("build signalr client: %w", err)
	}

	states := make(chan signalr.ClientState, 8)
	stopObserving := client.ObserveStateChanged(states)
	defer stopObserving()

	client.Start()
	defer client.Stop()

	for {
		select {
		case <-ctx.Done():
			s.connected.Store(false)
			return ctx.Err()
		case state := <-states:
			if state != signalr.ClientConnected {
				if s.connected.Swap(false) {
					log.Printf("Stream disconnected: %v\n", state)
				}
				continue
			}
			s.reconnects.Add(1)
			s.connected.Store(true)
			log.Printf("Stream connected, subscribing to %d chargers\n", len(s.chargers))
			// Subscriptions do not survive a reconnect, so they are replayed
			// every time the connection comes back up.
			s.subscribeAll(ctx, client)
		}
	}
}

func (s *Stream) connect(ctx context.Context) (signalr.Connection, error) {
	token, err := s.client.Token()
	if err != nil {
		return nil, fmt.Errorf("stream auth: %w", err)
	}
	return signalr.NewHTTPConnection(ctx, s.url,
		signalr.WithHTTPClient(&http.Client{Timeout: 60 * time.Second}),
		signalr.WithHTTPHeaders(func() http.Header {
			return http.Header{"Authorization": []string{"Bearer " + token.AccessToken}}
		}),
	)
}

func (s *Stream) subscribeAll(ctx context.Context, client signalr.Client) {
	for _, charger := range s.chargers {
		select {
		case <-ctx.Done():
			return
		case err := <-client.Send("SubscribeWithCurrentState", charger, true):
			if err != nil {
				log.Printf("Stream subscribe failed for %s: %v\n", charger, err)
			}
		case <-time.After(10 * time.Second):
			log.Printf("Stream subscribe timed out for %s\n", charger)
		}
	}
}

// handleProductUpdate applies one update. Easee sends a single observation
// per call, but an array is accepted too so a change in their framing cannot
// blind the exporter.
func (s *Stream) handleProductUpdate(raw json.RawMessage) {
	var single Observation
	if err := json.Unmarshal(raw, &single); err == nil && single.Mid != "" {
		s.applyObservation(&single)
		return
	}
	var batch []Observation
	if err := json.Unmarshal(raw, &batch); err == nil && len(batch) > 0 {
		for i := range batch {
			s.applyObservation(&batch[i])
		}
		return
	}
	log.Printf("Ignoring unrecognised product update: %s\n", raw)
}

func (s *Stream) applyObservation(o *Observation) {
	if o.Mid == "" {
		return
	}
	if o.Timestamp.IsZero() {
		o.Timestamp = time.Now()
	}
	s.apply(o.Mid, o)
}

// signalrLogger keeps the SignalR library quiet. Left to itself it debug
// logs every frame to stderr, which runs to hundreds of kilobytes per
// connection and buries the exporter's own output. Only entries carrying an
// error are passed through; connection state is reported by the stream
// metrics instead.
type signalrLogger struct{}

func (signalrLogger) Log(keyVals ...any) error {
	for i := 0; i+1 < len(keyVals); i += 2 {
		if fmt.Sprint(keyVals[i]) == "error" {
			log.Printf("Stream error: %v\n", keyVals[i+1])
			return nil
		}
	}
	return nil
}

// streamReceiver is called by the SignalR machinery. Method names must match
// the server's client-side calls exactly, and every method the server may
// call needs to exist: an unhandled call makes the library tear the
// connection down and reconnect, losing updates while it does.
type streamReceiver struct {
	signalr.Receiver
	stream *Stream
}

func (r *streamReceiver) ProductUpdate(raw json.RawMessage) {
	r.stream.handleProductUpdate(raw)
}

func (r *streamReceiver) ChargerUpdate(raw json.RawMessage) {
	r.stream.handleProductUpdate(raw)
}

func (r *streamReceiver) CommandResponse(raw json.RawMessage) {}

// SubscribeToMyProduct is called by the server and carries nothing the
// exporter needs, but it must be handled: leaving it out is what makes the
// connection drop and reconnect under load.
func (r *streamReceiver) SubscribeToMyProduct(raw json.RawMessage) {}
