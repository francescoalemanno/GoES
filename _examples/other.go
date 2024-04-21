package main

import (
	"fmt"

	GoES "github.com/francescoalemanno/GoES"
)

func abs2(x float64) float64 {
	return x * x
}
func myCustomFunction(x []float64) float64 {
	A, B := 1.0, -4.0
	return abs2(x[0]-A) + 100*abs2(x[1]+x[0]-A-B)
}

func main() {
	mu := []float64{1.0, 2.0}    // Initial mean vector
	sigma := []float64{0.5, 0.5} // Initial standard deviation vector

	optimizedMu, optimizedSigma := GoES.DefaultOpt(myCustomFunction, mu, sigma)

	fmt.Println("Optimized mean:", optimizedMu)
	fmt.Println("Optimized standard deviation:", optimizedSigma)
}
