package goes

import (
	"math"
	"testing"
)

func abs2(x float64) float64 {
	return x * x
}

func newFunction(ince float64, iva_detratta float64) float64 {
	ince_l := ince / (1 - 1.22/100)
	parte_incentivata := 0.65
	cost := func(x []float64) float64 {
		pc := x[0]
		fatt := pc + ince
		impo := fatt / iva_detratta
		pen := impo*parte_incentivata - ince_l
		return abs2(pen)
	}
	mu, _ := DefaultOpt(cost, []float64{2 * ince}, []float64{ince / 10})
	return mu[0]
}

func TestIncentivi(t *testing.T) {
	incentivi := []float64{1010.22, 1010.22, 1106.73, 2020.44}
	perc := []float64{1.0, 1.1, 1.22, 1.0}
	wants := []float64{563.15985, 720.49783, 996.17249, 1126.31970}
	for i := range incentivi {
		got := newFunction(incentivi[i], perc[i])
		want := wants[i]
		err := math.Abs(got - want)
		if err > 1e-5 {
			t.Errorf("got %.5f, wanted %.5f, err %.2g", got, want, err)
		}
	}
}

func TestBi(t *testing.T) {
	muw := []float64{4, -3}
	mu, sig := DefaultOpt(func(f []float64) float64 {
		return abs2(f[0]-muw[0]) + 100.0*abs2(f[0]+f[1]-muw[0]-muw[1])
	}, []float64{0.0, 0.0}, []float64{1.0, 1.0})
	err := math.Sqrt(abs2((mu[0]-muw[0])/muw[0]) + abs2((mu[1]-muw[1])/muw[1]))
	if err > 1e-6 {
		t.Error("got: ", mu, sig, " wanted:", muw, " error:", err)
	}
}
