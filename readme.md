**Package GoES - [![Go Report Card](https://goreportcard.com/badge/github.com/francescoalemanno/GoES)](https://goreportcard.com/report/github.com/francescoalemanno/GoES)**

**Overview**

The `GoES` package implements an Evolutionary Algorithm (EA) for optimization problems. It employs a diagonal variant of the Covariance Matrix Adaptation - Evolution Strategy (CMA-ES) algorithm, known for its efficiency and robustness in various optimization tasks.

**Functionality**

* **`Opt` function:** This core function performs the CMA-ES optimization. It takes the following arguments:
    * `fn`: A user-defined function representing the objective function to be optimized. This function should accept a slice of `float64` values as input and return a single `float64` value representing the cost or fitness of the solution.
    * `mu`: An initial mean vector of `float64` values, defining the starting point of the search in the solution space.
    * `sigma`: An initial standard deviation vector of `float64` values, determining the initial search radius around the mean vector.
    * `cfg` (optional): A configuration object of type `Config` (see below) to customize optimization parameters.
* **`DefaultOpt` function:** This convenience function calls `Opt` with default configuration values suitable for many common optimization problems.
   * `fn`: A user-defined function representing the objective function to be optimized. This function should accept a slice of `float64` values as input and return a single `float64` value representing the cost or fitness of the solution.
    * `mu`: An initial mean vector of `float64` values, defining the starting point of the search in the solution space.
    * `sigma`: An initial standard deviation vector of `float64` values, determining the initial search radius around the mean vector.
* **`TunedOpt` function:** This convenience function calls `Opt` with a configuration which is tuned according to the user cost-function. Running this tuned optimised is more expensive than using `Opt` directly but it may lead to better results.
   * `fn`: A user-defined function representing the objective function to be optimized. This function should accept a slice of `float64` values as input and return a single `float64` value representing the cost or fitness of the solution.
    * `mu`: An initial mean vector of `float64` values, defining the starting point of the search in the solution space.
    * `sigma`: An initial standard deviation vector of `float64` values, determining the initial search radius around the mean vector.
**Configuration**

The `Config` struct allows fine-tuning the optimization process:

* `Generations`: The maximum number of generations for the EA to run (default: 300).
* `PopSize`: The population size (number of candidate solutions) per generation (default: automatically determined based on problem dimensionality).
* `LR_mu`: Learning rate for mean vector update (default: 0.6).
* `LR_sigma`: Learning rate for standard deviation vector update (default: 0.15).
* `Momentum`: Momentum coefficient for velocity update (default: 0.93).
* `SigmaTol`: Tolerance threshold for stopping the optimization (default: 1e-12).
* `Verbose`: Flag to enable detailed logging of optimization progress during each generation (default: false).

**Usage Examples**

**Example 1: Sphere Function Optimization**

This example minimizes the sphere function, a common benchmark for optimization algorithms:

```go
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
	res, _ := GoES.Opt(sphere, mu, sigma, cfg)

	fmt.Println("Optimum:", res.Mu) // should be close to vector [0, 1, 2, ..., dim-1]
}
```

**Example 2: Default Optimization Config**

This example demonstrates how to use `GoES` to optimize a custom function:

```go
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

	res, _ := GoES.DefaultOpt(myCustomFunction, mu, sigma)

	fmt.Println("Optimized mean:", res.Mu)
	fmt.Println("Optimized standard deviation:", res.Sigma)
}
```
