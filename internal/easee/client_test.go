package easee

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"golang.org/x/oauth2"
)

func testClient() *Client {
	c := &Client{token: &oauth2.Token{AccessToken: "test", Expiry: time.Now().Add(time.Hour)}}
	c.HttpClient = oauth2.NewClient(context.Background(), c)
	return c
}

// A removed endpoint answers with an empty body, which used to surface as a
// bare "EOF" from the json decoder instead of the status code.
func TestDoRequestReportsEmptyErrorResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusGone)
	}))
	defer srv.Close()

	var out struct{}
	err := testClient().doRequest(srv.URL, nil, &out)
	if err == nil {
		t.Fatal("doRequest() = nil, want an error for 410 Gone")
	}
	if !strings.Contains(err.Error(), "410") {
		t.Errorf("doRequest() error = %q, want it to report status 410", err)
	}
}

// A rate limited request answers with a json object, which used to surface as
// "cannot unmarshal object into Go value of type []easee.Charger".
func TestDoRequestReportsErrorBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte(`{"title":"Too many requests"}`))
	}))
	defer srv.Close()

	var out []Charger
	err := testClient().doRequest(srv.URL, nil, &out)
	if err == nil {
		t.Fatal("doRequest() = nil, want an error for 429")
	}
	if !strings.Contains(err.Error(), "429") || !strings.Contains(err.Error(), "Too many requests") {
		t.Errorf("doRequest() error = %q, want status and body", err)
	}
}

func TestDoRequestDecodesSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`[{"id":"ECRTU44X","name":"Home","productCode":100}]`))
	}))
	defer srv.Close()

	var chargers []Charger
	if err := testClient().doRequest(srv.URL, nil, &chargers); err != nil {
		t.Fatalf("doRequest() = %v, want nil", err)
	}
	if len(chargers) != 1 || chargers[0].Id != "ECRTU44X" {
		t.Errorf("chargers = %+v, want one charger ECRTU44X", chargers)
	}
}

func TestChargerStateUrl(t *testing.T) {
	url := chargerStateUrl("ECRTU44X")
	if !strings.HasPrefix(url, "https://api.easee.com/state/ECRTU44X/observations?ids=") {
		t.Errorf("chargerStateUrl() = %q, unexpected host or path", url)
	}
	if !strings.HasSuffix(url, "?ids="+chargerStateObservationIds) {
		t.Errorf("chargerStateUrl() = %q, want all observation ids in one request", url)
	}
}
