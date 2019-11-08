package dlframework

// Copyright ©2014 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	math "math"
)

// JensenShannon computes the JensenShannon divergence between the distributions
// p and q. The Jensen-Shannon divergence is defined as
//  m = 0.5 * (p + q)
//  JS(p, q) = 0.5 ( KL(p, m) + KL(q, m) )
// Unlike Kullback-Liebler, the Jensen-Shannon distance is symmetric. The value
// is between 0 and ln(2).
func JensenShannon(p, q []float64) float64 {
	if len(p) != len(q) {
		panic("stat: slice length mismatch")
	}
	var js float64
	for i, v := range p {
		qi := q[i]
		m := 0.5 * (v + qi)
		if v != 0 {
			// add kl from p to m
			js += 0.5 * v * (math.Log(v) - math.Log(m))
		}
		if qi != 0 {
			// add kl from q to m
			js += 0.5 * qi * (math.Log(qi) - math.Log(m))
		}
	}
	return js
}

// KullbackLeibler computes the Kullback-Leibler distance between the
// distributions p and q. The natural logarithm is used.
//  sum_i(p_i * log(p_i / q_i))
// Note that the Kullback-Leibler distance is not symmetric;
// KullbackLeibler(p,q) != KullbackLeibler(q,p)
func KullbackLeibler(p, q []float64) float64 {
	if len(p) != len(q) {
		panic("stat: slice length mismatch")
	}
	var kl float64
	for i, v := range p {
		if v != 0 { // Entropy needs 0 * log(0) == 0
			kl += v * (math.Log(v) - math.Log(q[i]))
		}
	}
	return kl
}

// Covariance returns the weighted covariance between the samples of x and y.
//  sum_i {w_i (x_i - meanX) * (y_i - meanY)} / (sum_j {w_j} - 1)
// The lengths of x and y must be equal. If weights is nil then all of the
// weights are 1. If weights is not nil, then len(x) must equal len(weights).
func Covariance(x, y, weights []float64) float64 {
	// This is a two-pass corrected implementation.  It is an adaptation of the
	// algorithm used in the MeanVariance function, which applies a correction
	// to the typical two pass approach.

	if len(x) != len(y) {
		panic("stat: slice length mismatch")
	}
	xu := Mean(x, weights)
	yu := Mean(y, weights)
	var (
		ss            float64
		xcompensation float64
		ycompensation float64
	)
	if weights == nil {
		for i, xv := range x {
			yv := y[i]
			xd := xv - xu
			yd := yv - yu
			ss += xd * yd
			xcompensation += xd
			ycompensation += yd
		}
		// xcompensation and ycompensation are from Chan, et. al.
		// referenced in the MeanVariance function.  They are analogous
		// to the second term in (1.7) in that paper.
		return (ss - xcompensation*ycompensation/float64(len(x))) / float64(len(x)-1)
	}

	var sumWeights float64

	for i, xv := range x {
		w := weights[i]
		yv := y[i]
		wxd := w * (xv - xu)
		yd := (yv - yu)
		ss += wxd * yd
		xcompensation += wxd
		ycompensation += w * yd
		sumWeights += w
	}
	// xcompensation and ycompensation are from Chan, et. al.
	// referenced in the MeanVariance function.  They are analogous
	// to the second term in (1.7) in that paper, except they use
	// the sumWeights instead of the sample count.
	return (ss - xcompensation*ycompensation/sumWeights) / (sumWeights - 1)
}

// Mean computes the weighted mean of the data set.
//  sum_i {w_i * x_i} / sum_i {w_i}
// If weights is nil then all of the weights are 1. If weights is not nil, then
// len(x) must equal len(weights).
func Mean(x, weights []float64) float64 {
	if weights == nil {
		return Sum(x) / float64(len(x))
	}
	if len(x) != len(weights) {
		panic("stat: slice length mismatch")
	}
	var (
		sumValues  float64
		sumWeights float64
	)
	for i, w := range weights {
		sumValues += w * x[i]
		sumWeights += w
	}
	return sumValues / sumWeights
}

// bhattacharyyaCoeff computes the Bhattacharyya Coefficient for probability distributions given by:
//  \sum_i \sqrt{p_i q_i}
//
// It is assumed that p and q have equal length.
func bhattacharyyaCoeff(p, q []float64) float64 {
	var bc float64
	for i, a := range p {
		bc += math.Sqrt(a * q[i])
	}
	return bc
}

// Bhattacharyya computes the distance between the probability distributions p and q given by:
//  -\ln ( \sum_i \sqrt{p_i q_i} )
//
// The lengths of p and q must be equal. It is assumed that p and q sum to 1.
func Bhattacharyya(p, q []float64) float64 {
	if len(p) != len(q) {
		panic("stat: slice length mismatch")
	}
	bc := bhattacharyyaCoeff(p, q)
	return -math.Log(bc)
}

// Hellinger computes the distance between the probability distributions p and q given by:
//  \sqrt{ 1 - \sum_i \sqrt{p_i q_i} }
//
// The lengths of p and q must be equal. It is assumed that p and q sum to 1.
func Hellinger(p, q []float64) float64 {
	if len(p) != len(q) {
		panic("stat: slice length mismatch")
	}
	bc := bhattacharyyaCoeff(p, q)
	return math.Sqrt(1 - bc)
}

// Correlation returns the weighted correlation between the samples of x and y
// with the given means.
//  sum_i {w_i (x_i - meanX) * (y_i - meanY)} / (stdX * stdY)
// The lengths of x and y must be equal. If weights is nil then all of the
// weights are 1. If weights is not nil, then len(x) must equal len(weights).
func Correlation(x, y, weights []float64) float64 {
	// This is a two-pass corrected implementation.  It is an adaptation of the
	// algorithm used in the MeanVariance function, which applies a correction
	// to the typical two pass approach.

	if len(x) != len(y) {
		panic("stat: slice length mismatch")
	}
	xu := Mean(x, weights)
	yu := Mean(y, weights)
	var (
		sxx           float64
		syy           float64
		sxy           float64
		xcompensation float64
		ycompensation float64
	)
	if weights == nil {
		for i, xv := range x {
			yv := y[i]
			xd := xv - xu
			yd := yv - yu
			sxx += xd * xd
			syy += yd * yd
			sxy += xd * yd
			xcompensation += xd
			ycompensation += yd
		}
		// xcompensation and ycompensation are from Chan, et. al.
		// referenced in the MeanVariance function.  They are analogous
		// to the second term in (1.7) in that paper.
		sxx -= xcompensation * xcompensation / float64(len(x))
		syy -= ycompensation * ycompensation / float64(len(x))

		return (sxy - xcompensation*ycompensation/float64(len(x))) / math.Sqrt(sxx*syy)

	}

	var sumWeights float64
	for i, xv := range x {
		w := weights[i]
		yv := y[i]
		xd := xv - xu
		wxd := w * xd
		yd := yv - yu
		wyd := w * yd
		sxx += wxd * xd
		syy += wyd * yd
		sxy += wxd * yd
		xcompensation += wxd
		ycompensation += wyd
		sumWeights += w
	}
	// xcompensation and ycompensation are from Chan, et. al.
	// referenced in the MeanVariance function.  They are analogous
	// to the second term in (1.7) in that paper, except they use
	// the sumWeights instead of the sample count.
	sxx -= xcompensation * xcompensation / sumWeights
	syy -= ycompensation * ycompensation / sumWeights

	return (sxy - xcompensation*ycompensation/sumWeights) / math.Sqrt(sxx*syy)
}

// Sum returns the sum of the elements of the slice.
func Sum(s []float64) float64 {
	var sum float64
	for _, val := range s {
		sum += val
	}
	return sum
}
