package main

import (
	"fmt"
	"math/rand"

	GoES "github.com/francescoalemanno/GoES"
)

func abs2(x float64) float64 {
	return x * x
}
func myCustomFunction(x []float64) float64 {
	p := GoES.Probability(x[0])
	tot_runs := 110
	h := 0.0
	for _ = range tot_runs {
		if rand.Float64() < p {
			h += 1
		}
	}
	return abs2(h/float64(tot_runs) - 50.0/110.0)
}

func main() {
	mu := []float64{0}    // Initial mean vector
	sigma := []float64{1} // Initial standard deviation vector
	cfg := GoES.Config()
	cfg.Verbose = true
	cfg.Momentum = 0
	cfg.LR_mu = 1.0
	cfg.LR_sigma = 1.0
	cfg.PopSize = 200
	cfg.Generations = 200
	res, _ := GoES.Opt(myCustomFunction, mu, sigma, cfg)
	fmt.Println("Optimized mean:", GoES.Probability(res.Mu[0]))
	fmt.Println("MLE:", 50.0/110.0)
}
