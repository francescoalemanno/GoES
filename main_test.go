package GoES

import (
	"bytes"
	"log"
	"math"
	"strings"
	"testing"
)

func abs2(x float64) float64 {
	return x * x
}

func cost_test(ince float64, iva_detratta float64) float64 {
	ince_l := ince / (1 - 1.22/100)
	perc_pay := 0.65
	cost := func(x []float64) float64 {
		pc := x[0]
		fatt := pc + ince
		impo := fatt / iva_detratta
		pen := impo*perc_pay - ince_l
		return abs2(pen)
	}
	sol, err := Opt(cost, []float64{ince * 0.9}, []float64{ince / 10}, Config())
	if err != nil {
		return math.NaN()
	}
	return sol.Mu[0]
}

func TestUni(t *testing.T) {
	money_inc := []float64{1010.22, 1010.22, 1106.73, 2020.44}
	perc := []float64{1.0, 1.1, 1.22, 1.0}
	wants := []float64{563.15985, 720.49783, 996.17249, 1126.31970}
	for i := range money_inc {
		got := cost_test(money_inc[i], perc[i])
		want := wants[i]
		err := math.Abs(got - want)
		if err > 1e-5 {
			t.Errorf("got %.5f, wanted %.5f, err %.2g", got, want, err)
		}
	}
}
func TestBounded(t *testing.T) {
	cost := func(x []float64) float64 {
		f := Bounded(x[0], -2, 5)
		return f
	}
	res, _ := Opt(cost, []float64{0.0}, []float64{1.0}, Config())
	got := cost(res.Mu)
	want := -2.0
	err := math.Abs(got - want)
	if err > 1e-5 {
		t.Errorf("got %.5f, wanted %.5f, err %.2g", got, want, err)
	}
}

func TestBi(t *testing.T) {
	muw := []float64{4, -3}
	sol, err_opt := Opt(func(f []float64) float64 {
		return abs2(f[0]-muw[0]) + 100.0*abs2(f[0]+f[1]-muw[0]-muw[1])
	}, []float64{0.0, 0.0}, []float64{1.0, 1.0}, Config())
	if err_opt != nil {
		t.Error(err_opt)
	}
	mu := sol.Mu
	sig := sol.Sigma
	err := math.Sqrt(abs2((mu[0]-muw[0])/muw[0]) + abs2((mu[1]-muw[1])/muw[1]))
	if err > 1e-3 {
		t.Error("got: ", mu, sig, " wanted:", muw, " error:", err)
	}
}

func TestVerbose(t *testing.T) {
	cfg := Config()
	buf := bytes.NewBuffer([]byte{})
	log.SetOutput(buf)
	cfg.Verbose = true
	cfg.Generations = 10
	//Run it for few generations with a cost function to test for iteration verbosity
	Opt(func(f []float64) float64 {
		return Probability(f[0])
	}, []float64{0.0, 0.0}, []float64{1.0, 1.0}, cfg)
	//run it with a constant cost to hit convergence verbosity
	Opt(func(f []float64) float64 {
		return 0.0
	}, []float64{0.0, 0.0}, []float64{1.0, 1.0}, cfg)
	str := buf.String()
	if strings.Count(str, "GoES:") != cfg.Generations || strings.Count(str, "END OPT:") != 1 {
		t.Error(str)
	}
}

func TestError(t *testing.T) {
	_, err := Opt(
		func(f []float64) float64 { return 0.0 },
		[]float64{0.0, 0.0},
		[]float64{1.0}, // here there is a missing element
		Config(),
	)
	if err == nil {
		t.Error("Expected error")
	}
	_, err = Opt(
		func(f []float64) float64 { return 0.0 },
		[]float64{0.0, 0.0},
		[]float64{1.0, 1.0}, // here there is a missing element
		Config(),
	)
	if err != nil {
		t.Error("Expected success")
	}
}
