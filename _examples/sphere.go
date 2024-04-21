package main

import (
	"fmt"

	GoES "github.com/francescoalemanno/GoES"
)

func sphere(x []float64) float64 {
	var sum float64
	for int_i, v := range x {
		i := float64(int_i)
		sum += (v - i) * (v - i)
	}
	return sum
}

func main() {
	dim := 10                     // Dimensionality of the problem (number of variables)
	mu := make([]float64, dim)    // Initialize mean vector with zeros
	sigma := make([]float64, dim) // Initialize standard deviation vector with ones
	for i := range dim {
		sigma[i] = 1.0
	}

	// Optionally customize configuration
	cfg := GoES.Defaults()
	cfg.Generations = 1000
	cfg.Verbose = false

	// Perform optimization
	optimizedMu, _ := GoES.Opt(sphere, mu, sigma, cfg)

	fmt.Println("Optimum:", optimizedMu) // should be close to vector [0, 1, 2, ..., dim-1]
}
