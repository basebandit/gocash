package api

import "errors"

const (
	//current api version
	Version = "0.1"
)

type Money struct {
	Rates map[string]interface{}
	Base  string
}

func (m *Money) getRate(to, from string) (rate float64, fxErr error) {
	// Return an error if to rate isn't in the rates array
	if _, ok := m.Rates[to]; !ok {
		fxErr = errors.New("fx error: missing to rate")
		return
	}

	// Return an error if from rate isn't in the rates array
	if _, ok := m.Rates[from]; !ok {
		fxErr = errors.New("fx error: missing from rate")
		return
	}

	// If `from` currency == base, return the basic exchange rate for the `to` currency
	if from == m.Base {
		to := m.Rates[to]
		rate = to.(float64)
		return
	}

	// If `to` currency == fx.base, return the basic inverse rate of the `from` currency
	if to == m.Base {
		from := m.Rates[from]
		r := from.(float64)
		rate = 1 / r
		return
	}

	// Otherwise, return the `to` rate multipled by the inverse of the `from` rate to get the
	// relative exchange rate between the two currencies
	t := m.Rates[to]
	rTo := t.(float64)

	f := m.Rates[from]
	rFrm := f.(float64)

	rate = rTo * (1 / rFrm)
	return
}

//Convert converts a value from one currency to another
func (m *Money) Convert(val float64, from string, to string) (rate float64, cnvtErr error) {
	if from == "" && to == "" {
		cnvtErr = errors.New("fx error: missing either from or to rate")
		return
	}
	r, err := m.getRate(to, from)

	if err != nil {
		cnvtErr = err
		return
	}
	rate = val * r
	return
}
