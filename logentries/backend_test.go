package logentries_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/nu7hatch/gouuid"
	"github.com/pjvds/tidy"
	"github.com/pjvds/tidy/logentries"
)

func getLogTail(key string) (string, error) {
	url := fmt.Sprintf("https://pull.logentries.com/%v/hosts/ManualHost/tidy_tests/?start=-100000", key)
	response, err := http.Get(url)

	if err != nil {
		return "", err
	}

	all, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return "", err
	}

	return string(all), nil
}

func TestBackendRoundtrip(t *testing.T) {
	token := os.Getenv("LE_TOKEN")
	if len(token) == 0 {
		t.Skip("LE_TOKEN not set")
	}

	key := os.Getenv("LE_ACCOUNT_KEY")
	if len(key) == 0 {
		t.Skip("LE_ACCOUNT_KEY not set")
	}

	backend := logentries.Configure(token).UDP().Build()
	log := tidy.NewLogger("foobar", backend)

	id, _ := uuid.NewV4()
	secret := id.String()
	log.WithField("secret", secret).Debug("hello world")

	done := make(chan struct{})
	matched := make(chan struct{})
	go func() {
		t.Logf("looking for secret: %v", secret)
		for {
			select {
			case <-done:
				return
			default:
				tail, err := getLogTail(key)
				if err != nil {
					t.Fatalf("failed to get tail from logentries: %v", err)
				}

				if strings.Contains(tail, secret) {
					close(matched)
					return
				}

				time.Sleep(500 * time.Millisecond)
			}
		}
	}()

	select {
	case <-matched:
		// great
	case <-time.After(30 * time.Second):
		t.Fatalf("entry not found in tail")
		close(done)
	}
}
