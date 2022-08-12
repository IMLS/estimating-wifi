package api

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/spf13/viper"
	"gsa.gov/18f/cmd/session-counter/state"
)

func TestSimpleCall(t *testing.T) {
	durations := []*state.Duration{
		{
			SessionID: 0,
			Start:     1,
			End:       2,
		},
	}

	viper.Set("api.scheme", "https")
	viper.Set("api.host", "10x.gsa.gov")
	viper.Set("api.uri", "test/durations")

	var body string

	httpmock.Activate()
	httpmock.RegisterResponder("POST",
		"https://10x.gsa.gov/test/durations",
		func(req *http.Request) (*http.Response, error) {
			b, _ := io.ReadAll(req.Body)
			body = string(b)
			return httpmock.NewJsonResponse(200, `{}`)
		},
	)

	err := PostDurations(durations)
	if err != nil {
		t.Fatal("posting durations failed")
	}

	if !strings.Contains("start_time,end_time\n", body) {
		t.Fatal("posting durations did not contain csv header")
	}
	if !strings.Contains("1,2\n", body) {
		t.Fatal("posting durations did not contain csv body")
	}
}
