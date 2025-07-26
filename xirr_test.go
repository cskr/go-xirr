// Copyright 2018 Chandra Sekar S
// Use of this source code is governed by the license
// that can be found in the LICENSE file.

package xirr

import (
	"encoding/csv"
	"io"
	"math"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestSamples(t *testing.T) {
	cases := []struct {
		file string
		rate float64
	}{
		{"single_redemption.csv", 0.1361695793742},
		{"random.csv", 0.6924974337277},
		{"non_converging.csv", math.NaN()},
	}

	for _, c := range cases {
		t.Run(c.file, func(t *testing.T) {
			payments, err := loadPayments(c.file)
			if err != nil {
				t.Fatal("Error loading input:", err)
			}

			rate, err := Compute(payments)
			if err != nil {
				t.Fatal("Error computing XIRR:", err)
			}

			if math.IsNaN(c.rate) {
				if !math.IsNaN(rate) {
					t.Fatalf("Expected NaN, but was %.10f", rate)
				}
				return
			}

			if math.IsNaN(rate) || math.Abs(rate-c.rate) >= maxError {
				t.Fatalf("Expected %.10f, but was %.10f", c.rate, rate)
			}
		})
	}
}

func TestSameSign(t *testing.T) {
	_, err := Compute([]Payment{
		{parseDate("2016-06-11"), -100},
		{parseDate("2018-06-11"), -200},
	})
	if err != ErrInvalidPayments {
		t.Errorf("Invalid error for negative payments: %v", err)
	}

	_, err = Compute([]Payment{
		{parseDate("2016-06-11"), 100},
		{parseDate("2018-06-11"), 200},
	})
	if err != ErrInvalidPayments {
		t.Errorf("Invalid error for positive payments: %v", err)
	}
}

func loadPayments(file string) ([]Payment, error) {
	f, err := os.Open("samples/" + file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var payments []Payment
	r := csv.NewReader(f)
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		payments = append(payments, Payment{
			parseDate(rec[0]), parseAmount(rec[1]),
		})
	}

	return payments, nil
}

func parseDate(date string) time.Time {
	result, err := time.Parse("2006-01-02", date)
	if err != nil {
		panic(err)
	}
	return result
}

func parseAmount(num string) float64 {
	result, err := strconv.ParseFloat(num, 64)
	if err != nil {
		panic(err)
	}
	return result
}
