/*
This code implements a specific algorithm called CMA-ES (Covariance Matrix Adaptation Evolution Strategy) for optimization in Go. Here's a breakdown of what the code does:

**Configuration:**

* It defines a `Config` struct to hold various parameters for the optimization process.
* A `Defaults` function provides default values for these configurations.

**Function Optimization:**

* The core function is `Opt`. It takes several arguments:
  - `fn`: This is the function you want to optimize. It takes a list of floats as input and returns a single float representing the cost or fitness.
  - `mu`: The initial mean vector for the population.
  - `sigma`: The initial standard deviation vector for the population.
  - `cfg`: The configuration object with parameters like number of generations and learning rate.

* The function first calculates the population size based on the problem dimension.
* It defines helper functions for sorting and sampling candidate solutions.
* The main loop iterates for a specified number of generations:
  - It performs Nesterov momentum update on the mean vector.
  - It samples new candidate solutions based on the current mean and standard deviation.
  - It evaluates the cost of each candidate solution using the provided function `fn`.
  - It sorts the candidate solutions based on their cost (better solutions come first).
  - It updates the mean and standard deviation vectors based on the weighted contributions of promising candidate solutions.
  - It checks for convergence based on a minimum standard deviation threshold.
  - Optionally, it prints information about the current generation and best solution so far.

* Finally, the function returns the optimized mean (`mu`) and standard deviation (`sigma`) vectors.

**Convenience Function:**

* `DefaultOpt` is a wrapper function that uses `Defaults` to get a configuration and then calls `Opt` with some adjustments to the number of generations.

**Helper Functions:**

* `makeWeights` calculates weights for each candidate solution based on its rank in the population. This helps focus the search on promising regions.

Overall, this code provides an implementation of CMA-ES for optimizing a black-box function in Go. It allows you to specify your own objective function and configure various parameters for the optimization process.
*/
package goes

import (
	"fmt"
	"math"
	"slices"

	rand "golang.org/x/exp/rand"
)

type Config struct {
	Generations int
	PopSize     int
	LR_mu       float64
	LR_sigma    float64
	Momentum    float64
	SigmaTol    float64
	Verbose     bool
}

func Defaults() Config {
	cfg := Config{}
	cfg.Generations = 300
	cfg.PopSize = 0
	cfg.LR_mu = 0.6
	cfg.LR_sigma = 0.15
	cfg.Momentum = 0.93
	cfg.SigmaTol = 1e-12
	cfg.Verbose = false
	return cfg
}

const const_Ez0 = 0.7978845608028661 // mean(abs(randn()))
func Opt(fn func([]float64) float64, mu []float64, sigma []float64, cfg Config) ([]float64, []float64) {
	pop_n := cfg.PopSize
	n := len(mu)
	if len(sigma) != n {
		panic("mu and sigma must have the same length.")
	}
	for pop_n*pop_n <= 144*n {
		pop_n++
	}

	type Pair struct {
		Z []float64
		C float64
	}
	sortfn := func(a, b Pair) int {
		if a.C < b.C {
			return -1
		}
		if a.C > b.C {
			return 1
		}
		return 0
	}

	sample := func(av, sd []float64) Pair {
		z := make([]float64, n)
		trial := make([]float64, n)
		for {
			for i := range n {
				z[i] = rand.NormFloat64()
				trial[i] = z[i]*sd[i] + av[i]
			}
			cost := fn(trial)
			if !math.IsInf(cost, 0) && !math.IsNaN(cost) {
				return Pair{z, cost}
			}
		}
	}
	W := makeWeights(pop_n)
	pop := make([]Pair, pop_n)
	g := make([]float64, n)
	v := make([]float64, n)
	g_log_sigma := make([]float64, n)
	for runs := range cfg.Generations {
		nesterov_mu := make([]float64, n)
		for j := range n {
			v[j] *= cfg.Momentum
			nesterov_mu[j] = mu[j] + v[j]
		}
		for j := range pop_n {
			pop[j] = sample(nesterov_mu, sigma)
		}
		slices.SortFunc(pop, sortfn)

		for j := range n {
			g[j] = 0
			g_log_sigma[j] = 0
		}
		for i, p := range pop {
			if W[i] < 0 {
				break
			}
			for j := range n {
				g[j] += W[i] * p.Z[j]
				g_log_sigma[j] += W[i] * (math.Abs(p.Z[j])/const_Ez0 - 1)
			}
		}
		for j := range n {
			v[j] += cfg.LR_mu * sigma[j] * g[j]
			mu[j] += v[j]
			sigma[j] *= math.Exp(cfg.LR_sigma * g_log_sigma[j])
		}
		if slices.Max(sigma) < cfg.SigmaTol {
			break
		}
		if cfg.Verbose {
			fmt.Println(runs, mu, sigma, pop[pop_n/2].C)
		}
	}
	return mu, sigma
}

func DefaultOpt(fn func([]float64) float64, mu []float64, sigma []float64) ([]float64, []float64) {
	cfg := Defaults()
	cfg.Generations = int(math.Ceil(math.Sqrt(float64(len(mu)*2+1)) * 300))
	return Opt(fn, mu, sigma, cfg)
}
func makeWeights(pop_size int) []float64 {
	W := make([]float64, pop_size)
	for i := range pop_size {
		W[i] = math.Log(float64(pop_size-1)*0.5+0.5) - math.Log(float64(i)+0.5)
		if W[i] < 0 {
			W[i] = 0
		}
	}
	sumW := float64(0)
	for _, v := range W {
		sumW += v
	}
	for i := range pop_size {
		W[i] /= sumW
	}
	return W
}
