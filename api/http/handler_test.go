package http

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"testing"

	"github.com/ubombi/timeseries/api"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
)

var event = `{
  "event_type": "failure",
  "ts": 7774465743,
  "params": {
  	"": "+380976658510",
	"username": "Vitalii",
	"password": "qwerty",
	"location": "Ukraine",
	"timezone": "Europe/Kiev"

  }
}`

var malformed = `{
  "event_type": "failure",
  "ts": 7774465743,
  "params":, {
  	+380976658510",
	"timezone": "Europe/Kiev"

  }
}`

type testStorage struct {
	called int
	err    error
}

func (ts *testStorage) Store(e api.Event) error {
	ts.called += 1
	return ts.err
}

func TestService_HandlerFunc(t *testing.T) {
	tests := []struct {
		name        string
		storage     testStorage
		req         *http.Request
		storeCalled bool
		statusCode  int
	}{
		{"normal event", testStorage{}, skipErr(http.NewRequest("POST", "http://127.0.0.1/", strings.NewReader(event))), true, 200},
		{"empty json", testStorage{}, skipErr(http.NewRequest("POST", "http://127.0.0.1/", strings.NewReader("{}"))), false, 400},
		{"malformed json", testStorage{}, skipErr(http.NewRequest("POST", "http://127.0.0.1/", strings.NewReader(malformed))), false, 400},
		{"malformed json", testStorage{}, skipErr(http.NewRequest("POST", "http://127.0.0.1/", nil)), false, 400},
		{"malformed json", testStorage{err: errors.New("some error")}, skipErr(http.NewRequest("POST", "http://127.0.0.1/", strings.NewReader(event))), true, 500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				Storage: &tt.storage,
			}
			resp, err := serve(s.HandlerFunc, tt.req)
			if err != nil {
				t.Error(err)
			}
			if resp.StatusCode != tt.statusCode {
				t.Error("wrong status code. ", resp.StatusCode, tt.statusCode)
			}
			if tt.storeCalled && tt.storage.called == 0 {
				t.Error("item wasnt stored")
			}
			if !tt.storeCalled && tt.storage.called != 0 {
				t.Error("item should not be stored")
			}
		})
	}
}

// serve serves http request using provided fasthttp handler
func serve(handler fasthttp.RequestHandler, req *http.Request) (*http.Response, error) {
	ln := fasthttputil.NewInmemoryListener()
	defer ln.Close()

	go func() {
		err := fasthttp.Serve(ln, handler)
		if err != nil {
			panic(fmt.Errorf("failed to serve: %v", err))
		}
	}()

	client := http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return ln.Dial()
			},
		},
	}

	return client.Do(req)
}

func skipErr(r *http.Request, err error) *http.Request {
	return r
}
