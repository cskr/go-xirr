// Copyright 2018 Chandra Sekar S
// Use of this source code is governed by the license
// that can be found in the LICENSE file.

package xirr

import (
	"errors"
	"math"
	"sort"
	"time"
)

const maxError = 1e-10

var ErrInvalidPayments = errors.New("negative and positive payments are required")

type Payment struct {
	Date   time.Time
	Amount float64
}

func Compute(payments []Payment) (xirr float64, err error) {
	if err := validatePayments(payments); err != nil {
		return 0, err
	}

	rate := computeWithGuess(payments, 0.1)
	for guess := -0.99; guess < 1.0 && (math.IsNaN(rate) || math.IsInf(rate, 0)); guess += 0.1 {
		rate = computeWithGuess(payments, guess)
	}

	return rate, nil
}

func validatePayments(payments []Payment) error {
	positive, negative := false, false
	for _, p := range payments {
		if p.Amount >= 0.0 {
			positive = true
		}
		if p.Amount < 0.0 {
			negative = true
		}
	}

	if !positive || !negative {
		return ErrInvalidPayments
	}
	return nil
}

func computeWithGuess(payments []Payment, guess float64) float64 {
	sorted := make([]Payment, len(payments))
	copy(sorted, payments)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Date.Before(sorted[j].Date)
	})

	r, e := guess, 1.0
	for e > maxError {
		r1 := r - xirr(sorted, r)/dxirr(sorted, r)
		e = math.Abs(r1 - r)
		r = r1
	}

	return r
}

func xirr(payments []Payment, rate float64) float64 {
	result := 0.0
	for _, p := range payments {
		exp := getExp(p, payments[0])
		result += p.Amount / math.Pow(1.0+rate, exp)
	}
	return result
}

func dxirr(payments []Payment, rate float64) float64 {
	result := 0.0
	for _, p := range payments {
		exp := getExp(p, payments[0])
		result -= p.Amount * exp / math.Pow(1.0+rate, exp+1.0)
	}
	return result
}

func getExp(p, p0 Payment) float64 {
	return float64(p.Date.Sub(p0.Date)/(24*time.Hour)) / 365
}
