package load

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	vegeta "github.com/tsenart/vegeta/lib"
)

func TestHealth(t *testing.T) {
	rate := vegeta.Rate{Freq: 50, Per: time.Second}
	duration := 30 * time.Second
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: http.MethodGet,
		URL:    "http://localhost:8080/api/v1/health",
	})

	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics

	// Create a file to save all results
	//f, err := os.Create("results.bin")
	//if err != nil {
	//	panic(err)
	//}
	//defer f.Close()

	//encoder := vegeta.NewEncoder(f)

	for res := range attacker.Attack(targeter, rate, duration, "Shorten API Test") {
		metrics.Add(res)
		//if err := encoder.Encode(res); err != nil {
		//	fmt.Printf("Failed to encode result: %v\n", err)
		//}
	}
	metrics.Close()

	// Print summary
	fmt.Printf("Requests: %d\n", metrics.Requests)
	fmt.Printf("Success rate: %.2f%%\n", metrics.Success*100)
	fmt.Printf("Avg latency: %s\n", metrics.Latencies.Mean)
	fmt.Printf("P99 latency: %s\n", metrics.Latencies.P99)
}
