package main

import (
	"encoding/json"
	"flag"
	"log"
	"math/rand/v2"
	"net/http"
)

type MultiplierService struct {
	rtp float64
}

type Response struct {
	Result float64 `json:"result"`
}

func main() {
	rtp := flag.Float64("rtp", 0.95, "RTP value (0 < rtp <= 1.0)")
	flag.Parse()
	if *rtp <= 0.0 || *rtp > 1.0 {
		log.Fatalf("0 < rtp <= 1.0")
	}

	service := NewMultiplierService(*rtp)

	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		multiplier := service.GenerateMultiplier()

		response := Response{Result: multiplier}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
		}
	})

	log.Printf("Server is listening...")
	log.Fatal(http.ListenAndServe(":64333", nil))
}

func NewMultiplierService(rtp float64) *MultiplierService {
	return &MultiplierService{rtp: rtp}
}

func (m *MultiplierService) GenerateMultiplier() float64 {
	u := rand.Float64()
	if u >= m.rtp {
		return 1.0
	}

	var multiplier float64

	for {
		uPareto := rand.Float64()
		if uPareto == 0.0 {
			uPareto = 1e-10
		}

		multiplier = (1.0 / uPareto) * 2.0
		if multiplier <= 10000.0 {
			break
		}
	}

	if multiplier < 1.0 {
		multiplier = 1.0
	}

	return multiplier
}
