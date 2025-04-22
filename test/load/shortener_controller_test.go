package load

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	vegeta "github.com/tsenart/vegeta/lib"
)

func TestShortenerController_ShortURL(t *testing.T) {
	rate := vegeta.Rate{Freq: 50, Per: time.Second}
	duration := 30 * time.Second

	targeter := func(tgt *vegeta.Target) error {
		tgt.Method = http.MethodPost
		tgt.URL = "http://localhost:8080/api/v1/shorten"
		tgt.Body = []byte(fmt.Sprintf(`{"longUrl": "%s"}`, gofakeit.URL()))
		return nil
	}

	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics

	for res := range attacker.Attack(targeter, rate, duration, "Shorten API Test") {
		metrics.Add(res)
	}
	metrics.Close()

	// Print summary
	fmt.Printf("Requests: %d\n", metrics.Requests)
	fmt.Printf("Success rate: %.2f%%\n", metrics.Success*100)
	fmt.Printf("Avg latency: %s\n", metrics.Latencies.Mean)
	fmt.Printf("P99 latency: %s\n", metrics.Latencies.P99)
}
