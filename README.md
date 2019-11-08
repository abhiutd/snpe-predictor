# go-snpe
Go bindings for Qualcomm SNPE C++ API. This code is being used to perform ML inference on mobile devices.

# Installation

## TfLite

## Dependencies

## Bindings

# User API

New() // create a new predictor with appropriate ML model and hardware mode
Predict() // run prediction
ReadPredictionOutput() // return top-1 prediction
Close() // close predictor

Refer to cbits.go for API signatures.

# ML inference
